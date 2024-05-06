package makesdk

import (
	"fmt"
	"time"
)

type ScenarioLog struct {
	Id             string    `json:"id"`
	ImtId          string    `json:"imtId"`
	Duration       int       `json:"duration"`
	Operations     int       `json:"operations"`
	Transfer       int       `json:"transfer"`
	OrganizationId int       `json:"organizationId"`
	TeamId         int       `json:"teamId"`
	Type           string    `json:"type"`
	AuthorId       int       `json:"authorId"`
	Instant        bool      `json:"instant"`
	Timestamp      time.Time `json:"timestamp"`
	Status         int       `json:"status"`
	ScenarioId     int
}

type ScenarioLogResponse struct {
	ScenarioLog ScenarioLog `json:"scenarioLog"`
}

type ScenarioLogListResponse struct {
	ScenarioLogs []ScenarioLog `json:"scenarioLogs"`
	Pagination   Pagination    `json:"pg"`
}

type ScenarioLogListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *ScenarioLogListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *ScenarioLogListPaginator) NextPage() ([]ScenarioLog, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ScenarioLogListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeScenariosRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.ScenarioLogs)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.ScenarioLogs, nil
}

func (at *Client) NewScenarioLogListPaginator(maxItems int, scenarioId int) *ScenarioLogListPaginator {
	var config = NewRequestConfig(fmt.Sprintf(`scenarios/%d/logs`, scenarioId))

	maxItems, _ = GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(50)

	var p = &ScenarioLogListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetScenarioLog(scenarioId int, executionId string) (*ScenarioLog, error) {
	var config = NewRequestConfig(fmt.Sprintf(`scenarios/%d/logs/%s`, scenarioId, executionId))

	var result = &ScenarioLogResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, TokenScopeScenariosRead)
	}

	return &result.ScenarioLog, nil
}
