package errors

import (
	"errors"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customError"
	"hyperledger-api-skeleton/hyperledger-business-logic/utils/golang/customLogger"
)

const (
	InternalError              customError.ErrorCode = 1000
	TimeoutError                                     = 1001
	FileOperationError                               = 1002
	DataManipulationError                            = 1003
	ValidationError                                  = 1004
	AuthorisationError                               = 1005
	HttpRequestError                                 = 1006
	BlockchainError                                  = 1007
	KafkaError                                       = 1008
	PostgresError                                    = 1009
	PostgresDuplicatedKeyError                       = 1010
	PostgresRecordNotFound                           = 1011
	RedisError                                       = 1012
	RedisKeyNotFound                                 = 1013
)

const (
	extServerError customError.ErrorCode = 2001
)

var externalErrorMessage = map[customError.ErrorCode]string{
	extServerError: "Internal Server Error",
}

var internalAndExternalErrorList = map[customError.ErrorCode]customError.Translator{
	InternalError:              {InternalMessage: "internal error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	TimeoutError:               {InternalMessage: "timeout error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	FileOperationError:         {InternalMessage: "file operation error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	DataManipulationError:      {InternalMessage: "data manipulation error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	ValidationError:            {InternalMessage: "validation error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	AuthorisationError:         {InternalMessage: "authorisation error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	HttpRequestError:           {InternalMessage: "http request error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	BlockchainError:            {InternalMessage: "blockchain error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	KafkaError:                 {InternalMessage: "kafka error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	PostgresError:              {InternalMessage: "postgres error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	PostgresDuplicatedKeyError: {InternalMessage: "duplicate key error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	PostgresRecordNotFound:     {InternalMessage: "record not found error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	RedisError:                 {InternalMessage: "redis error", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	RedisKeyNotFound:           {InternalMessage: "redis key not found", ExternalError: customError.ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
}

func NewCustomError(code customError.ErrorCode, description string) customError.Error {
	return customError.NewError(code, description,
		customError.WithErrorList(internalAndExternalErrorList),
		customError.WithCallerLevel(customError.CallerLevel3))
}

func HandlerPostgresErr(err error) customError.Error {
	errC := NewCustomError(PostgresError, err.Error())

	var pgError *pgconn.PgError
	switch {
	case errors.As(err, &pgError):
		if errors.Is(err, pgError) {
			switch pgError.Code {
			case "23505":
				errC = NewCustomError(PostgresDuplicatedKeyError, err.Error())
			default:
			}
		}
	default:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errC = NewCustomError(PostgresRecordNotFound, err.Error())
		}
	}
	return errC
}

func GatewayHandleError(ctx context.Context, log customLogger.Log, err error) customError.Error {
	switch err := err.(type) {
	case *client.EndorseError:
		log.ErrorCtx(ctx, "Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err.Error())
		return NewCustomError(EndorseError, err.Error())
	case *client.SubmitError:
		log.ErrorCtx(ctx, "Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err.Error())
		return NewCustomError(SubmitError, err.Error())
	case *client.CommitStatusError:
		if goErrors.Is(err, context.DeadlineExceeded) {
			log.ErrorCtx(ctx, "Timeout waiting for transaction %s commit status: %s", err.TransactionID, err.Error())
		} else {
			log.ErrorCtx(ctx, "Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err.Error())
		}
		return NewCustomError(CommitStatusError, err.Error())
	case *client.CommitError:
		log.ErrorCtx(ctx, "Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err.Error())
		return NewCustomError(CommitError, err.Error())
	default:
		log.ErrorCtx(ctx, "unexpected error type %T: %w", err, err.Error())
		return NewCustomError(BlockchainError, err.Error())
	}
}
