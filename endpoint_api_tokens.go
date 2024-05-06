package makesdk

import (
	"fmt"
)

type ApiToken struct {
	Token    string       `json:"token"`
	Label    string       `json:"label"`
	Scope    []TokenScope `json:"scope"`
	Created  string       `json:"created"`
	IsActive bool
}

type TokenScope string

const (
	TokenScopeAdminRead                     TokenScope = "admin:read"
	TokenScopeAdminWrite                    TokenScope = "admin:write"
	TokenScopeAgentsRead                    TokenScope = "agents:read"
	TokenScopeAgentsWrite                   TokenScope = "agents:write"
	TokenScopeAppsRead                      TokenScope = "apps:read"
	TokenScopeAppsWrite                     TokenScope = "apps:write"
	TokenScopeConnectionsRead               TokenScope = "connections:read"
	TokenScopeConnectionsWrite              TokenScope = "connections:write"
	TokenScopeCustomPropertyStructuresRead  TokenScope = "custom-property-structures:read"
	TokenScopeCustomPropertyStructuresWrite TokenScope = "custom-property-structures:write"
	TokenScopeDataStoresRead                TokenScope = "datastores:read"
	TokenScopeDataStoresWrite               TokenScope = "datastores:write"
	TokenScopeDevicesRead                   TokenScope = "devices:read"
	TokenScopeDevicesWrite                  TokenScope = "devices:write"
	TokenScopeDlqsRead                      TokenScope = "dlqs:read"
	TokenScopeDlqsWrite                     TokenScope = "dlqs:write"
	TokenScopeFunctionsRead                 TokenScope = "functions:read"
	TokenScopeFunctionsWrite                TokenScope = "functions:write"
	TokenScopeHooksRead                     TokenScope = "hooks:read"
	TokenScopeHooksWrite                    TokenScope = "hooks:write"
	TokenScopeImtFormsRead                  TokenScope = "imt-forms:read"
	TokenScopeInstancesRead                 TokenScope = "instances:read"
	TokenScopeInstancesWrite                TokenScope = "instances:write"
	TokenScopeKeysRead                      TokenScope = "keys:read"
	TokenScopeKeysWrite                     TokenScope = "keys:write"
	TokenScopeNotificationsRead             TokenScope = "notifications:read"
	TokenScopeNotificationsWrite            TokenScope = "notifications:write"
	TokenScopeOrganizationsRead             TokenScope = "organizations:read"
	TokenScopeOrganizationsWrite            TokenScope = "organizations:write"
	TokenScopeOrganizationVariablesRead     TokenScope = "organization-variables:read"
	TokenScopeOrganizationVariablesWrite    TokenScope = "organization-variables:write"
	TokenScopeScenariosRead                 TokenScope = "scenarios:read"
	TokenScopeScenariosRun                  TokenScope = "scenarios:run"
	TokenScopeScenariosWrite                TokenScope = "scenarios:write"
	TokenScopeSdkAppsRead                   TokenScope = "sdk-apps:read"
	TokenScopeSdkAppsWrite                  TokenScope = "sdk-apps:write"
	TokenScopeSystemRead                    TokenScope = "system:read"
	TokenScopeSystemWrite                   TokenScope = "system:write"
	TokenScopeTeamsRead                     TokenScope = "teams:read"
	TokenScopeTeamsWrite                    TokenScope = "teams:write"
	TokenScopeTeamVariablesRead             TokenScope = "team-variables:read"
	TokenScopeTeamVariablesWrite            TokenScope = "team-variables:write"
	TokenScopeTemplatesRead                 TokenScope = "templates:read"
	TokenScopeTemplatesWrite                TokenScope = "templates:write"
	TokenScopeUdtsRead                      TokenScope = "udts:read"
	TokenScopeUdtsWrite                     TokenScope = "udts:write"
	TokenScopeUserRead                      TokenScope = "user:read"
	TokenScopeUserWrite                     TokenScope = "user:write"
)

type ApiTokenListResponse struct {
	ApiTokens []ApiToken `json:"apiTokens"`
}

type ApiTokenListPaginator struct {
	firstPage bool
	config    *RequestConfig
	client    *Client
}

func (lp *ApiTokenListPaginator) HasMorePages() bool {
	return lp.firstPage
}

func (lp *ApiTokenListPaginator) NextPage() ([]ApiToken, error) {
	if !lp.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	var r = &ApiTokenListResponse{}
	var _, err = lp.client.Get(lp.config, r)
	if err != nil {
		return nil, lp.client.handleKnownErrors(err, TokenScopeUserRead)
	}

	lp.firstPage = false

	return r.ApiTokens, nil
}

func (at *Client) NewApiTokenListPaginator(_ int) *ApiTokenListPaginator {
	var config = NewRequestConfig("users/me/api-tokens")

	var p = &ApiTokenListPaginator{
		firstPage: true,
		config:    config,
		client:    at,
	}

	return p
}
