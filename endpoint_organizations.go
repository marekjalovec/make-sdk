package makesdk

import "fmt"

type Organization struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	CountryId   int                 `json:"countryId"`
	TimezoneId  int                 `json:"timezoneId"`
	License     OrganizationLicence `json:"license"`
	Zone        string              `json:"zone"`
	ServiceName string              `json:"serviceName"`
	IsPaused    bool                `json:"isPaused"`
	ExternalId  string              `json:"externalId"`
	Teams       []OrganizationTeam  `json:"teams"` // used to load dependant objects
}

type OrganizationTeam struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type OrganizationLicence struct {
	ApiLimit           int64  `json:"apiLimit"`
	CreatingTemplates  bool   `json:"creatingTemplates"`
	CustomFunctions    bool   `json:"customFunctions"`
	CustomVariables    bool   `json:"customVariables"`
	DsLimit            int64  `json:"dslimit"`
	DssLimit           int64  `json:"dsslimit"`
	ExecutionTime      int    `json:"executionTime"`
	FsLimit            int64  `json:"fslimit"`
	Fulltext           bool   `json:"fulltext"`
	GracePeriod        int    `json:"gracePeriod"`
	InstallPublicApps  bool   `json:"installPublicApps"`
	Interval           int    `json:"interval"`
	IoLimit            int64  `json:"iolimit"`
	OnDemandScheduling bool   `json:"onDemandScheduling"`
	Operations         int64  `json:"operations"`
	RestartPeriod      string `json:"restartPeriod"`
	ScenarioIO         bool   `json:"scenarioIO"`
	Transfer           int64  `json:"transfer"`
}

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Pagination    Pagination     `json:"pg"`
}

type OrganizationListPaginator struct {
	firstPage  bool
	maxItems   int
	totalCount int
	lastCount  int
	config     *RequestConfig
	client     *Client
}

func (lp *OrganizationListPaginator) HasMorePages() bool {
	var fullPage = lp.lastCount == lp.config.Pagination.Limit
	var allLoaded = lp.totalCount == lp.maxItems

	return lp.firstPage || (fullPage && !allLoaded)
}

func (lp *OrganizationListPaginator) NextPage() ([]Organization, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &OrganizationListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeOrganizationsRead)
	}

	lp.firstPage = false
	lp.lastCount = len(r.Organizations)
	var newPageSize = -1
	if lp.totalCount+lp.config.Pagination.Limit > lp.maxItems {
		newPageSize = lp.maxItems - lp.totalCount
	}
	lp.config.Pagination.NextPage(newPageSize)
	lp.totalCount += lp.lastCount

	return r.Organizations, nil
}

func (at *Client) NewOrganizationListPaginator(maxItems int) *OrganizationListPaginator {
	var config = NewRequestConfig("organizations")
	ColumnsToParams(config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})

	maxItems, limit := GetMaxAndLimit(maxItems)
	config.Pagination = NewRequestPagination(limit)

	var p = &OrganizationListPaginator{
		firstPage:  true,
		maxItems:   maxItems,
		totalCount: 0,
		lastCount:  0,
		config:     config,
		client:     at,
	}

	return p
}

func (at *Client) GetOrganization(organizationId int) (*Organization, error) {
	var config = NewRequestConfig(fmt.Sprintf(`organizations/%d`, organizationId))
	ColumnsToParams(config.Params, []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId", "teams"})

	var result = &OrganizationResponse{}
	var _, err = at.Get(config, &result)
	if err != nil {
		return nil, at.handleKnownErrors(err, TokenScopeOrganizationsRead)
	}

	return &result.Organization, nil
}
