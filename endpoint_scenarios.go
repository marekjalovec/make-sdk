package makesdk

import (
	"fmt"
	"strconv"
	"time"
)

type Scenario struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	TeamId       int       `json:"teamId"`
	HookId       int       `json:"hookId,omitempty"`
	DeviceId     int       `json:"deviceId,omitempty"`
	DeviceScope  string    `json:"deviceScope,omitempty"`
	Description  string    `json:"description"`
	FolderId     int       `json:"folderId,omitempty"`
	IsInvalid    bool      `json:"isinvalid"`
	IsLinked     bool      `json:"islinked"`
	IsLocked     bool      `json:"islocked"`
	IsPaused     bool      `json:"isPaused"`
	Concept      bool      `json:"concept"`
	UsedPackages []string  `json:"usedPackages"`
	LastEdit     time.Time `json:"lastEdit"`
	Scheduling   struct {
		Type     string `json:"type"`
		Interval int    `json:"interval"`
	} `json:"scheduling"`
	IsWaiting     bool `json:"iswaiting"`
	DlqCount      int  `json:"dlqCount"`
	CreatedByUser struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"createdByUser"`
	UpdatedByUser struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"updatedByUser"`
	NextExec         time.Time `json:"nextExec"`
	ScenarioVersion  int       `json:"scenarioVersion"`
	ModuleSequenceId int       `json:"moduleSequenceId"`
	OrganizationId   int
}

type ScenarioResponse struct {
	Scenario Scenario `json:"scenario"`
}

type ScenarioListResponse struct {
	Scenarios  []Scenario `json:"scenarios"`
	Pagination Pagination `json:"pg"`
}

type ScenarioListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *ScenarioListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *ScenarioListPaginator) NextPage() ([]Scenario, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ScenarioListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "scenarios:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.Scenarios)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Scenarios, nil
}

func (at *Client) NewScenarioListPaginator(maxItems int, teamId int, organizationId int) *ScenarioListPaginator {
	var config = NewRequestConfig("scenarios")
	if organizationId != 0 {
		config.Params.Set("organizationId", strconv.Itoa(organizationId))
	} else if teamId != 0 {
		config.Params.Set("teamId", strconv.Itoa(teamId))
	} else {
		return nil
	}

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &ScenarioListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetScenario(scenarioId int) (*Scenario, error) {
	var config = NewRequestConfig(fmt.Sprintf(`scenarios/%d`, scenarioId))

	var result = &ScenarioResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, "scenarios:read")
	}

	return &result.Scenario, nil
}
