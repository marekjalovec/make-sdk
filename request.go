package makesdk

import (
	"net/url"
)

type RequestConfig struct {
	Endpoint   string
	Params     url.Values
	Pagination *RequestPagination
}

func NewRequestConfig(endpoint string) *RequestConfig {
	return &RequestConfig{
		Endpoint: endpoint,
		Params:   url.Values{},
	}
}

type RequestPagination struct {
	Limit  int
	Offset int
}

func (rp *RequestPagination) NextPage(newPageSize int) {
	rp.Offset += rp.Limit

	if newPageSize != -1 {
		rp.Limit = newPageSize
	}
}

func NewRequestPagination(limit int) *RequestPagination {
	return &RequestPagination{
		Limit:  limit,
		Offset: 0,
	}
}

type Pagination struct {
	SortBy  string `json:"sortBy"`
	Limit   int    `json:"limit"`
	SortDir string `json:"sortDir"`
	Offset  int    `json:"offset"`
}
