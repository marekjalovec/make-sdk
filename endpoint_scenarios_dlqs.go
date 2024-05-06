package makesdk

import (
	"fmt"
	"strconv"
	"time"
)

type ScenarioDlq struct {
	Id           string    `json:"id"`
	Reason       string    `json:"reason"`
	Size         int       `json:"size"`
	Index        int       `json:"index,omitempty"`
	Retry        bool      `json:"retry"`
	Attempts     int       `json:"attempts"`
	Created      time.Time `json:"created"`
	Resolved     bool      `json:"resolved"`
	Deleted      bool      `json:"deleted,omitempty"`
	ExecutionId  string    `json:"executionId,omitempty"`
	ScenarioId   int       `json:"scenarioId,omitempty"`
	ScenarioName string    `json:"scenarioName,omitempty"`
	TeamId       int       `json:"companyId,omitempty"`
	TeamName     string    `json:"companyName,omitempty"`
}

type ScenarioDlqResponse struct {
	ScenarioDlq ScenarioDlq `json:"dlq"`
}

type ScenarioDlqListResponse struct {
	ScenarioDlqs []ScenarioDlq `json:"dlqs"`
	Pagination   Pagination    `json:"pg"`
}

type ScenarioDlqListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *ScenarioDlqListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *ScenarioDlqListPaginator) NextPage() ([]ScenarioDlq, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ScenarioDlqListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeDlqsRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.ScenarioDlqs)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.ScenarioDlqs, nil
}

func (at *Client) NewScenarioDlqListPaginator(maxItems int, scenarioId int) *ScenarioDlqListPaginator {
	var config = NewRequestConfig("dlqs")
	config.Params.Set("scenarioId", strconv.Itoa(scenarioId))

	maxItems, _ = GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(50)

	var p = &ScenarioDlqListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetScenarioDlq(dlqId string) (*ScenarioDlq, error) {
	var config = NewRequestConfig(fmt.Sprintf(`dlqs/%s`, dlqId))

	var result = &ScenarioDlqResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, TokenScopeDlqsRead)
	}

	return &result.ScenarioDlq, nil
}
