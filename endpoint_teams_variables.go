package makesdk

import (
	"fmt"
)

type TeamVariable struct {
	Name     string `json:"name"`
	TypeId   int    `json:"typeId"`
	Value    any    `json:"value"`
	IsSystem bool   `json:"isSystem"`
	TeamId   int
}

type TeamVariableListResponse struct {
	TeamVariables []TeamVariable `json:"teamVariables"`
}

type TeamVariableListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *TeamVariableListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *TeamVariableListPaginator) NextPage() ([]TeamVariable, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &TeamVariableListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeTeamVariablesRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.TeamVariables)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.TeamVariables, nil
}

func (at *Client) NewTeamVariableListPaginator(maxItems int, teamId int) *TeamVariableListPaginator {
	var config = NewRequestConfig(fmt.Sprintf(`teams/%d/variables`, teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &TeamVariableListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
