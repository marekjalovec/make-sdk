package makesdk

import (
	"fmt"
)

type ApiToken struct {
	Token    string   `json:"token"`
	Label    string   `json:"label"`
	Scope    []string `json:"scope"`
	Created  string   `json:"created"`
	IsActive bool
}

type ApiTokenListResponse struct {
	ApiTokens []ApiToken `json:"apiTokens"`
}

type ApiTokenListPaginator struct {
	firstPage bool
	config    *RequestConfig
	client    *Client
}

func (lp *ApiTokenListPaginator) HasMorePages() bool {
	return lp.firstPage
}

func (lp *ApiTokenListPaginator) NextPage() ([]ApiToken, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ApiTokenListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "user:read")
	}

	lp.firstPage = false

	return r.ApiTokens, nil
}

func (at *Client) NewApiTokenListPaginator(_ int) *ApiTokenListPaginator {
	var config = NewRequestConfig("users/me/api-tokens")

	var p = &ApiTokenListPaginator{
		firstPage: true,
		config:    config,
		client:    at,
	}

	return p
}
