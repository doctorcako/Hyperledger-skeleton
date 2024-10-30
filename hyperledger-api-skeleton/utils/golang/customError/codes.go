package customError

type Translator struct {
	InternalMessage string
	ExternalError   ExtError
}

type ExtError struct {
	Code    ErrorCode
	Message string
}

const (
	InternalError         ErrorCode = 1000
	TimeoutError                    = 1001
	FileOperationError              = 1002
	DataManipulationError           = 1003
	AuthorisationError              = 1004
	HttpRequestError                = 1005
	DataBaseError                   = 1006
)

const (
	extServerError ErrorCode = 2001
)

var externalErrorMessage = map[ErrorCode]string{
	extServerError: "Internal Server Error",
}

var internalAndExternalErrorList = map[ErrorCode]Translator{
	InternalError:         {InternalMessage: "internal error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	TimeoutError:          {InternalMessage: "timeout error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	FileOperationError:    {InternalMessage: "file operation error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	DataManipulationError: {InternalMessage: "data manipulation error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	AuthorisationError:    {InternalMessage: "authorisation error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	HttpRequestError:      {InternalMessage: "http request error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
	DataBaseError:         {InternalMessage: "database error", ExternalError: ExtError{Code: extServerError, Message: externalErrorMessage[extServerError]}},
}
