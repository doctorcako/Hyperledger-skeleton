package customError

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type ErrorCode int
type CallerLevel int

type Error interface {
	Error() string
	ErrorCode() ErrorCode
	ErrorDescription() string
	ToError() error
	ExternalErrorCode() int
	ExternalErrorMsg() string
}

type errorCustom struct {
	code        ErrorCode
	message     string
	description string
	fileLine    string
	errorsList  map[ErrorCode]Translator
}

var (
	defaultErrorListOnce sync.Once
	defaultErrorList     map[ErrorCode]Translator
)

const (
	CallerLevel2 CallerLevel = 2
	CallerLevel3 CallerLevel = 3
)

type errorsOptions struct {
	errorsList  map[ErrorCode]Translator
	callerLevel CallerLevel
}

type OptionFunc func(*errorsOptions)

func defaultOpts() errorsOptions {

	//Initialises the default error messages
	defaultErrorListOnce.Do(initializeDefaultErrorMessages)

	errorOpts := errorsOptions{
		errorsList:  defaultErrorList,
		callerLevel: CallerLevel2,
	}

	return errorOpts
}

func WithErrorList(errorList map[ErrorCode]Translator) OptionFunc {
	return func(o *errorsOptions) {
		o.errorsList = errorList
	}
}

func WithCallerLevel(callerLevel CallerLevel) OptionFunc {
	return func(o *errorsOptions) {
		o.callerLevel = callerLevel
	}
}

// initializeDefaultErrorMessages default message
func initializeDefaultErrorMessages() {
	defaultErrorList = internalAndExternalErrorList
}

func NewError(code ErrorCode, description string, opts ...OptionFunc) Error {
	//default configuration
	o := defaultOpts()

	//Apply custom errors.
	for _, opt := range opts {
		opt(&o)
	}

	return newError(code, o.errorsList, description, o.callerLevel)
}

func newError(code ErrorCode, errorsList map[ErrorCode]Translator, description string, callerLevel CallerLevel) Error {

	//Get error line
	_, file, line, _ := runtime.Caller(int(callerLevel))
	pcs := make([]uintptr, 100)
	_ = runtime.Callers(3, pcs)

	//Get error message
	msg := errorsList[code]

	return &errorCustom{
		code:        code,
		message:     msg.InternalMessage,
		description: description,
		fileLine:    file + ":" + strconv.Itoa(line),
		errorsList:  errorsList,
	}
}

func (e *errorCustom) Error() string {
	return e.generateMsgString()
}

func (e *errorCustom) ErrorCode() ErrorCode {
	return e.code
}

func (e *errorCustom) ErrorDescription() string {
	return e.description
}

func (e *errorCustom) ToError() error {
	return errors.New(e.generateMsgString())
}

func (e *errorCustom) ExternalErrorCode() int {
	externalError, _ := e.errorsList[e.code]
	return int(externalError.ExternalError.Code)
}

func (e *errorCustom) ExternalErrorMsg() string {
	externalError, _ := e.errorsList[e.code]
	return externalError.ExternalError.Message
}

func (e *errorCustom) generateMsgString() string {
	var buf bytes.Buffer

	// Print the error code if there is one.
	if e.code != 0 {
		codeString := strconv.Itoa(int(e.code))
		buf.WriteString("<" + codeString + "> ")
	}

	// Print the file-line, if any.
	if e.fileLine != "" {
		buf.WriteString(e.fileLine + " - ")
	}

	// Print the message, if any.
	if e.message != "" {
		buf.WriteString(e.message + ": ")
	}

	// Print the original error message, if any.
	if e.description != "" {
		buf.WriteString(e.description)
	}

	return strings.TrimSpace(buf.String())
	//return strings.TrimSuffix(strings.TrimSpace(buf.String()), ",")
}
