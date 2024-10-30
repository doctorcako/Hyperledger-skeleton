package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Ok(body any)
	Created(body any)
	OkPagination(body any, page, limit, total int)
	Redirect301(body any)
}

// respInfo is the HTTP response information
type respInfo struct {
	Writer  http.ResponseWriter
	Headers map[string]string
}

// NewResponse creates and returns a new response

// HandlerExternalError convert internal to external error based on error catalog function

// WriteResponse - write HTTP headers and body
func (resp *respInfo) WriteResponse(code int, body interface{}, pag *Pagination) error {
	resp.writeHeaders()
	resp.writeStatusCode(code)

	respBody := responseHttp{}
	if body != nil {
		respBody.Body = body
	}
	if pag != nil {
		respBody.Pagination = pag
	}

	respBodyBytes, err := json.Marshal(respBody)

	return nil
}

func (resp *respInfo) writeHeaders() {
	for key, value := range resp.Headers {
		resp.Writer.Header().Set(key, value)
	}
}

func (resp *respInfo) writeStatusCode(code int) {
	resp.Writer.WriteHeader(code)
}
