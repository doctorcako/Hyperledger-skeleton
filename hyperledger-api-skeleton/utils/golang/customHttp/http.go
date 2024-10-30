package customHttp

import (
	"net/http"
)

const defaultMaxIdempotentCallAttempts int = 0

type opts struct {
	timeout           int
	maxRetries        int
	conditionsToRetry RetryIfFunc
	urlToLogRemotely  map[string]string
}

type RetryIfFunc func(request *http.Request) bool
type OptsFunc func(*opts)

func NewHttpProvider(opts ...OptsFunc) {
	o, errC := nil
	if errC != nil {
		return nil, errC
	}
	for _, fn := range opts {
		fn(&o)
	}
	return &o, nil
}

func WithLogRemotely(urls map[string]string) OptsFunc {
	return func(opts *opts) {
		opts.urlToLogRemotely = urls
	}
}

func WithTimeout(timeout int) OptsFunc {
	return func(opts *opts) {
		opts.timeout = timeout
	}
}

func WithRetries(totalRetries int) OptsFunc {
	return func(opts *opts) {
		opts.maxRetries = totalRetries
	}
}
