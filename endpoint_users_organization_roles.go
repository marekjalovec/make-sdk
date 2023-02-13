package makesdk

import (
	"fmt"
)

type UserOrganizationRole struct {
	UserId         int    `json:"userId"`
	UsersRoleId    int    `json:"usersRoleId"`
	OrganizationId int    `json:"organizationId"`
	Invitation     string `json:"invitation"`
	SsoPending     bool   `json:"ssoPending"`
}

type UserOrganizationRoleListResponse struct {
	UserOrganizationRoles []UserOrganizationRole `json:"userOrganizationRoles"`
	Pagination            Pagination             `json:"pg"`
}

type UserOrganizationRoleListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *UserOrganizationRoleListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *UserOrganizationRoleListPaginator) NextPage() ([]UserOrganizationRole, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &UserOrganizationRoleListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "user:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.UserOrganizationRoles)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.UserOrganizationRoles, nil
}

func (at *Client) NewUserOrganizationRoleListPaginator(maxItems int, userId int) *UserOrganizationRoleListPaginator {
	var config = NewRequestConfig(fmt.Sprintf(`users/%d/user-organization-roles`, userId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &UserOrganizationRoleListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
