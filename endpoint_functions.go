package makesdk

import (
	"fmt"
	"strconv"
	"time"
)

type Function struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Args        string `json:"args"`
	Scenarios   []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"scenarios"`
	CreatedByUser struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"createdByUser"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	TeamId    int
}

type FunctionResponse struct {
	Function Function `json:"function"`
}

type FunctionListResponse struct {
	Functions []Function `json:"functions"`
	//Pagination Pagination `json:"pg"` // not used now, API does not support pagination yet
}

type FunctionListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *FunctionListPaginator) HasMorePages() bool {
	return lp.firstPage
}

func (lp *FunctionListPaginator) NextPage() ([]Function, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &FunctionListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeFunctionsRead)
	}

	lp.firstPage = false

	return r.Functions, nil
}

func (at *Client) NewFunctionListPaginator(maxItems int, teamId int) *FunctionListPaginator {
	var config = NewRequestConfig("functions")
	config.Params.Set("teamId", strconv.Itoa(teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &FunctionListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetFunction(functionId int) (*Function, error) {
	var config = NewRequestConfig(fmt.Sprintf(`functions/%d`, functionId))

	var result = &FunctionResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, TokenScopeFunctionsRead)
	}

	return &result.Function, nil
}
