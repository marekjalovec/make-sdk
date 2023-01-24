package makesdk

import (
	"fmt"
	"strconv"
)

type DataStore struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Records         any    `json:"records"`
	Size            string `json:"size"`
	MaxSize         string `json:"maxSize"`
	DatastructureId int    `json:"datastructureId"`
	TeamId          int    `json:"teamId"`
}

type DataStoreResponse struct {
	DataStore DataStore `json:"dataStore"`
}

type DataStoreListResponse struct {
	DataStores []DataStore `json:"dataStores"`
	Pg         Pagination  `json:"pg"`
}

type DataStoreListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *DataStoreListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *DataStoreListPaginator) NextPage() ([]DataStore, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &DataStoreListResponse{}
	var err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "datastores:read")
	}

	lp.firstPage = false
	lp.lastCount = len(r.DataStores)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.DataStores, nil
}

func (at *Client) NewDataStoreListPaginator(maxItems int, teamId int) *DataStoreListPaginator {
	var config = NewRequestConfig("data-stores")
	config.Params.Set("teamId", strconv.Itoa(teamId))

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &DataStoreListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetDataStore(dataStoreId int) (*DataStore, error) {
	var config = NewRequestConfig(fmt.Sprintf(`data-stores/%d`, dataStoreId))
	ColumnsToParams(&config.Params, []string{"id", "name", "teamId", "records", "size", "maxSize", "datastructureId"})

	var result = &DataStoreResponse{}
	var err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, "datastores:read")
	}

	return &result.DataStore, nil
}
