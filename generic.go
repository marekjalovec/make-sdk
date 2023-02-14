package makesdk

import (
	"encoding/json"
	"fmt"
)

type GenericListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
	dst        func() GenericListResponse
}

func (lp *GenericListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *GenericListPaginator) NextPage() (interface{}, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var b, err = lp.client.Get(lp.config, nil)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, "unknown")
	}

	var pg = struct {
		Pagination Pagination `json:"pg"`
	}{}
	err = json.Unmarshal(b, &pg)
	if err != nil {
		return nil, fmt.Errorf("JSON decode failed: %s error: %w", b, err)
	}

	var r = lp.dst()
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, fmt.Errorf("JSON decode failed: %s error: %w", b, err)
	}

	lp.firstPage = false
	lp.lastCount = r.GetItemCount()
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r, nil
}

func (at *Client) NewGenericListPaginator(maxItems int, config *RequestConfig, dst func() GenericListResponse) *GenericListPaginator {
	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &GenericListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
		dst:        dst,
	}

	return p
}

type GenericListResponse interface {
	GetItemCount() int
}
