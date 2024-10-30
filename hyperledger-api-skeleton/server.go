package main

import (
	"context"
	appHttp "hyperledger-api-skeleton/hyperledger-business-logic/apps/http"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customError"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customLogger"

	"github.com/gin-gonic/gin"

	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/errors"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/properties"
	"hyperledger-api-skeleton/hyperledger-business-logic/pkg/sdk/fabricGateway"
	"os"
	"os/signal"
	"syscall"
)

const (
	ApplicationProperties string = "./configs/application_properties.yaml"
)

func NewBinder(properties *properties.Properties) appHttp.Binder {
	ctx := context.Background()

	//Define logger
	log := customLogger.NewLog(
		customLogger.LogModuleName(properties.ApplicationProperties.Micro.Name),
		customLogger.LogLevel(customLogger.StringToLogLevel(properties.ApplicationProperties.HttpServer.LogLevel)))

	optsGateway := fabricGateway.Opts{
		MspID: properties.ApplicationProperties.GatewayParams.MspID,
		// CryptoPath:   properties.ApplicationProperties.GatewayParams.CryptoPath,
		// CertPath:     properties.ApplicationProperties.GatewayParams.CertPath,
		// KeyPath:      properties.ApplicationProperties.GatewayParams.KeyPath,
		// TlsCertPath:  properties.ApplicationProperties.GatewayParams.TlsCertPath,
		// PeerEndpoint: properties.ApplicationProperties.GatewayParams.PeerEndpoint,
		GatewayPeer: properties.ApplicationProperties.GatewayParams.GatewayPeer,
		ChannelName: properties.ApplicationProperties.GatewayParams.ChannelName,
	}

	//Hyperledger Fabric Gateway
	fabricG, errC := fabricGateway.NewFabricGateway(ctx, log, optsGateway)
	if errC != nil {
		panic(errC.Error())
	}

	//Status service
	statusService := status.NewStatusService(log, properties, fabricG)

	// Kafka

	router := gin.New()

	return appHttp.Binder{
		Cxt:           ctx,
		Properties:    properties,
		LogCustom:     log,
		ServerHttp:    appHttp.ServerHttp{Router: router},
		FabricGateway: fabricG,
		StatusService: statusService,
	}
}

func main() {
	propertiesInfo := new(properties.Properties)
	errC := propertiesInfo.InitAllConfig(ApplicationProperties)
	if errC != nil {
		panic(errC.Error())
	}

	binder := NewBinder(propertiesInfo)

	//start listen network connections and communications

	networkChanges := make(chan os.Signal)
	signal.Notify(networkChanges, syscall.SIGINT, syscall.SIGTERM)
	binder.KafkaInput.StartListenEvents(networkChanges)

	defer binder.FabricGateway.CloseConnection()

	errs := make(chan customError.Error, 2)
	go func() {
		errs <- appHttp.StartApp(binder)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		err := <-c
		errs <- errors.NewCustomError(errors.InternalError, err.String())
	}()

	errResp := <-errs
	binder.LogCustom.Error("microservice shutdown " + errResp.Error())
}
