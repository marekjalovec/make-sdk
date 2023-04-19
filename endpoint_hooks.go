package makesdk

import (
	"fmt"
	"strconv"
)

type Hook struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	TeamId      int    `json:"teamId"`
	Udid        string `json:"udid"`
	Type        string `json:"type"`
	TypeName    string `json:"typeName"`
	PackageName string `json:"packageName"`
	Theme       string `json:"theme"`
	Flags       struct {
		Form bool `json:"form"`
	} `json:"flags"`
	IsEditable bool `json:"editable"`
	IsEnabled  bool `json:"enabled"`
	IsGone     bool `json:"gone"`
	QueueCount int  `json:"queueCount"`
	QueueLimit int  `json:"queueLimit"`
	Data       struct {
		Headers   bool   `json:"headers"`
		Method    bool   `json:"method"`
		Stringify bool   `json:"stringify"`
		TeamId    int    `json:"teamId"`
		Ip        string `json:"ip"`
		Udt       int    `json:"udt"`
	} `json:"data"`
	ScenarioId int    `json:"scenarioId"`
	Url        string `json:"url"`
}

type HookResponse struct {
	Hook Hook `json:"hook"`
}

type HookListResponse struct {
	Hooks      []Hook     `json:"hooks"`
	Pagination Pagination `json:"pg"`
}

type HookListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *HookListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *HookListPaginator) NextPage() ([]Hook, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &HookListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "hooks:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.Hooks)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Hooks, nil
}

func (at *Client) NewHookListPaginator(maxItems int, teamId int) *HookListPaginator {
	var config = NewRequestConfig("hooks")
	config.Params.Set("teamId", strconv.Itoa(teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &HookListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetHook(hookId int) (*Hook, error) {
	var config = NewRequestConfig(fmt.Sprintf(`hooks/%d`, hookId))

	var result = &HookResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, "hooks:read")
	}

	return &result.Hook, nil
}
