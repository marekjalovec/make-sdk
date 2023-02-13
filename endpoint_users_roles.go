package makesdk

import (
	"fmt"
)

type UserRole struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Subsidiary  bool             `json:"subsidiary"`
	Category    UserRoleCategory `json:"category"`
	Permissions []string         `json:"permissions"`
}

type UserRoleCategory string

const (
	UserRoleCategoryTeam         UserRoleCategory = "team"
	UserRoleCategoryOrganization UserRoleCategory = "organization"
)

type UserRoleListResponse struct {
	UserRoles  []UserRole `json:"usersRoles"`
	Pagination Pagination `json:"pg"`
}

type UserRoleListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *UserRoleListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *UserRoleListPaginator) NextPage() ([]UserRole, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &UserRoleListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "user:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.UserRoles)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.UserRoles, nil
}

func (at *Client) NewUserRoleListPaginator(maxItems int) *UserRoleListPaginator {
	var config = NewRequestConfig("users/roles")
	ColumnsToParams(config.Params, []string{"id", "name", "subsidiary", "category", "permissions"})

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &UserRoleListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
