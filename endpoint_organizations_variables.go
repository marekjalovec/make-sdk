package makesdk

import (
	"fmt"
)

type OrganizationVariable struct {
	Name           string `json:"name"`
	TypeId         int    `json:"typeId"`
	Value          any    `json:"value"`
	IsSystem       bool   `json:"isSystem"`
	OrganizationId int
}

type OrganizationVariableListResponse struct {
	OrganizationVariables []OrganizationVariable `json:"organizationVariables"`
}

type OrganizationVariableListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *OrganizationVariableListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *OrganizationVariableListPaginator) NextPage() ([]OrganizationVariable, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &OrganizationVariableListResponse{}
	var err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "organization-variables:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.OrganizationVariables)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.OrganizationVariables, nil
}

func (at *Client) NewOrganizationVariableListPaginator(maxItems int, organizationId int) *OrganizationVariableListPaginator {
	var config = NewRequestConfig(fmt.Sprintf(`organizations/%d/variables`, organizationId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &OrganizationVariableListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
