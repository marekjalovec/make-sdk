package makesdk

import (
	"fmt"
)

type UserTeamRole struct {
	UserId      int  `json:"userId"`
	UsersRoleId int  `json:"usersRoleId"`
	TeamId      int  `json:"teamId"`
	Changeable  bool `json:"changeable"`
}

type UserTeamRoleListResponse struct {
	UserTeamRoles []UserTeamRole `json:"userTeamRoles"`
	Pagination    Pagination     `json:"pg"`
}

type UserTeamRoleListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *UserTeamRoleListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *UserTeamRoleListPaginator) NextPage() ([]UserTeamRole, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &UserTeamRoleListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeUserRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.UserTeamRoles)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.UserTeamRoles, nil
}

func (at *Client) NewUserTeamRoleListPaginator(maxItems int, userId int) *UserTeamRoleListPaginator {
	var config = NewRequestConfig(fmt.Sprintf(`users/%d/user-team-roles`, userId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &UserTeamRoleListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
