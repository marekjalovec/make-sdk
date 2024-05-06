package makesdk

import (
	"fmt"
	"strconv"
)

type Team struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationId int    `json:"organizationId"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

type TeamListResponse struct {
	Teams      []Team     `json:"teams"`
	Pagination Pagination `json:"pg"`
}

type TeamListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *TeamListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *TeamListPaginator) NextPage() ([]Team, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &TeamListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeTeamsRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.Teams)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Teams, nil
}

func (at *Client) NewTeamListPaginator(maxItems int, organizationId int) *TeamListPaginator {
	var config = NewRequestConfig("teams")
	ColumnsToParams(config.Params, []string{"id", "name", "organizationId"})
	config.Params.Set("organizationId", strconv.Itoa(organizationId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &TeamListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetTeam(teamId int) (*Team, error) {
	var config = NewRequestConfig(fmt.Sprintf(`teams/%d`, teamId))
	ColumnsToParams(config.Params, []string{"id", "name", "organizationId"})

	var result = &TeamResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, TokenScopeTeamsRead)
	}

	return &result.Team, nil
}
