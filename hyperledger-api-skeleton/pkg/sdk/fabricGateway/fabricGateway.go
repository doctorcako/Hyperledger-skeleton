package fabricGateway

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/errors"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customError"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customLogger"
	"os"
	"path"
	"time"
)

type Opts struct {
	MspID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TlsCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	ChannelName  string
}

type gatewayFabric struct {
	mspID            string
	cryptoPath       string
	certPath         string
	keyPath          string
	tlsCertPath      string
	peerEndpoint     string
	gatewayPeer      string
	channelName      string
	clientConnection *grpc.ClientConn
	gwConnection     *client.Gateway
	network          *client.Network
	log              customLogger.Log
}

type GatewayFabric interface {
	ExecuteTx(ctx context.Context, chaincodeName string, functionName string, txArgs [][]byte) ([]byte, customError.Error)
	ExecuteQuery(ctx context.Context, chaincodeName string, functionName string, txArgs [][]byte) ([]byte, customError.Error)
	CloseConnection()
}

func NewFabricGateway(ctx context.Context, log customLogger.Log, opts Opts) (GatewayFabric, customError.Error) {
	headerLog := "NewFabricGateway"
	clientConnection, err := newGrpcConnection(ctx, opts, log)
	if err != nil {
		log.ErrorCtx(ctx, headerLog, "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	id, err := newIdentity(ctx, opts, log)
	if err != nil {
		log.ErrorCtx(ctx, headerLog, "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	sign, err := newSign(ctx, opts, log)
	if err != nil {
		log.ErrorCtx(ctx, headerLog, "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	// Create a Gateway connection for a specific client identity
	gw, errC := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if errC != nil {
		log.ErrorCtx(ctx, headerLog, "Error:", errC.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	network := gw.GetNetwork(opts.ChannelName)

	g := gatewayFabric{
		mspID:            opts.MspID,
		cryptoPath:       opts.CryptoPath,
		certPath:         opts.CertPath,
		keyPath:          opts.KeyPath,
		tlsCertPath:      opts.TlsCertPath,
		peerEndpoint:     opts.PeerEndpoint,
		gatewayPeer:      opts.GatewayPeer,
		channelName:      opts.ChannelName,
		clientConnection: clientConnection,
		gwConnection:     gw,
		network:          network,
		log:              log,
	}

	return &g, nil
}

func (f *gatewayFabric) ExecuteTx(ctx context.Context, chaincodeName string, functionName string, txArgs [][]byte) ([]byte, customError.Error) {
	headerLog := "gatewayFabric.ExecuteTx - GATEWAY - "
	f.log.DebugCtx(ctx, headerLog, "Init")

	// Convert [][]byte to []string
	args := make([]string, len(txArgs))
	for i, b := range txArgs {
		args[i] = string(b)
	}

	contract := f.network.GetContract(chaincodeName)
	response, err := contract.SubmitTransaction(functionName, args...)
	if err != nil {
		f.log.ErrorCtx(ctx, headerLog, "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, f.log, err)
	}

	f.log.DebugCtx(ctx, headerLog, "End")
	return response, nil
}

func (f *gatewayFabric) ExecuteQuery(ctx context.Context, chaincodeName string, functionName string, txArgs [][]byte) ([]byte, customError.Error) {
	headerLog := "gatewayFabric.ExecuteQuery - GATEWAY - "
	f.log.DebugCtx(ctx, headerLog, "Init")

	// Convert [][]byte to []string
	args := make([]string, len(txArgs))
	for i, b := range txArgs {
		args[i] = string(b)
	}

	contract := f.network.GetContract(chaincodeName)
	response, err := contract.EvaluateTransaction(functionName, args...)
	if err != nil {
		f.log.ErrorCtx(ctx, headerLog, "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, f.log, err)
	}

	f.log.DebugCtx(ctx, headerLog, "End")
	return response, nil
}

func (f *gatewayFabric) CloseConnection() {
	f.clientConnection.Close()
	f.gwConnection.Close()
}

func newGrpcConnection(ctx context.Context, opts Opts, log customLogger.Log) (*grpc.ClientConn, customError.Error) {
	certificate, err := loadCertificate(opts.TlsCertPath)
	if err != nil {
		log.ErrorCtx(ctx, "newGrpcConnection", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, opts.GatewayPeer)

	connection, err := grpc.Dial(opts.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		log.ErrorCtx(ctx, "newGrpcConnection - failed to create gRPC connection", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	return connection, nil
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(ctx context.Context, opts Opts, log customLogger.Log) (*identity.X509Identity, customError.Error) {
	certificate, err := loadCertificate(opts.CertPath)
	if err != nil {
		log.ErrorCtx(ctx, "newIdentity - error loading certificate", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	id, err := identity.NewX509Identity(opts.MspID, certificate)
	if err != nil {
		log.ErrorCtx(ctx, "newIdentity - error creating x509 identity", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	return id, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(ctx context.Context, opts Opts, log customLogger.Log) (identity.Sign, customError.Error) {
	files, err := os.ReadDir(opts.KeyPath)
	if err != nil {
		log.ErrorCtx(ctx, "newSign - failed to read private key directory", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}
	privateKeyPEM, err := os.ReadFile(path.Join(opts.KeyPath, files[0].Name()))

	if err != nil {
		log.ErrorCtx(ctx, "newSign - failed to read private key file", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.ErrorCtx(ctx, "newSign - failed to read private key from PEM", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		log.ErrorCtx(ctx, "newSign - failed to create private key sign", "Error:", err.Error())
		return nil, errors.GatewayHandleError(ctx, log, err)
	}

	return sign, nil
}
