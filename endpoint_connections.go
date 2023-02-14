package makesdk

import (
	"fmt"
	"strconv"
	"time"
)

type Connection struct {
	Id           int                `json:"id"`
	Name         string             `json:"name"`
	AccountName  string             `json:"accountName"`
	AccountLabel string             `json:"accountLabel"`
	PackageName  string             `json:"packageName"`
	Expire       time.Time          `json:"expire"`
	Metadata     ConnectionMetadata `json:"metadata,omitempty"`
	TeamId       int                `json:"teamId"`
	Upgradeable  bool               `json:"upgradeable"`
	Scoped       bool               `json:"scoped"`
	Scopes       []ConnectionScope  `json:"scopes,omitempty"`
	AccountType  string             `json:"accountType"`
	Editable     bool               `json:"editable"`
	Uid          string             `json:"uid"`
}

type ConnectionMetadata struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type ConnectionScope struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Account string `json:"account,omitempty"`
}

type ConnectionResponse struct {
	Connection Connection `json:"connection"`
}

type ConnectionListResponse struct {
	Connections []Connection `json:"connections"`
	Pg          Pagination   `json:"pg"`
}

type ConnectionListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *ConnectionListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *ConnectionListPaginator) NextPage() ([]Connection, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ConnectionListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "connections:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.Connections)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Connections, nil
}

func (at *Client) NewConnectionListPaginator(maxItems int, teamId int) *ConnectionListPaginator {
	var config = NewRequestConfig("connections")
	ColumnsToParams(config.Params, []string{"id", "name", "accountName", "accountLabel", "packageName", "expire", "metadata", "teamId", "upgradeable", "scoped", "accountType", "editable", "uid"})
	config.Params.Set("teamId", strconv.Itoa(teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &ConnectionListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetConnection(connectionId int) (*Connection, error) {
	var config = NewRequestConfig(fmt.Sprintf(`connections/%d`, connectionId))

	var result = &ConnectionResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, "connections:read")
	}

	return &result.Connection, nil
}
