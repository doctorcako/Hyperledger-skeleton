package response

type responseHttp struct {
	Body       interface{} `json:"body,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type errorCustom struct {
	Code    int
	Details errorDetails
}

type errorDetails struct {
	InternalCode   int
	InternalReason string
}

type RespHttp struct {
	Body       []byte `json:"body"`
	StatusCode int    `json:"statusCode"`
}

type Pagination struct {
	Page  int `json:"page"`
	Total int `json:"total"`
}

type RespBodyError struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details DetailsError `json:"details"`
}

type DetailsError struct {
	InternalCode   int    `json:"internalCode"`
	InternalReason string `json:"internalReason"`
	Tracker        string `json:"tracker"`
}
