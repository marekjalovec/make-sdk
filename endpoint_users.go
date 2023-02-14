package makesdk

import (
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Id             int          `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Language       string       `json:"language"`
	TimezoneId     int          `json:"timezoneId"`
	LocaleId       int          `json:"localeId"`
	CountryId      int          `json:"countryId"`
	Features       UserFeatures `json:"features"`
	Avatar         string       `json:"avatar"`
	LastLogin      time.Time    `json:"lastLogin"`
	OrganizationId int
	TeamId         int
}

type UserFeatures struct {
	AllowApps       bool `json:"allow_apps"`
	AllowAppsJs     bool `json:"allow_apps_js"`
	PrivateModules  bool `json:"private_modules"`
	AllowAppsCommit bool `json:"allow_apps_commit"`
	LocalAccess     bool `json:"local_access"`
}

type UserListResponse struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pg"`
}

type UserListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *UserListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *UserListPaginator) NextPage() ([]User, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &UserListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "user:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.Users)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Users, nil
}

func (at *Client) NewUserListPaginator(maxItems int, teamId int) *UserListPaginator {
	var config = NewRequestConfig("users")
	config.Params.Set("teamId", strconv.Itoa(teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &UserListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}
