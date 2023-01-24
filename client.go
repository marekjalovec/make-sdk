package makesdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultPageSize  = 10000
	defaultRateLimit = 50
)

type Client struct {
	client      *http.Client
	rateLimiter <-chan time.Time
	baseUrl     string
	apiToken    string
	scopes      *[]string
}

var clientInstance *Client

func GetClient(config *Config) *Client {
	if clientInstance != nil {
		return clientInstance
	}

	if config.RateLimit == nil {
		config.RateLimit = &defaultRateLimit
	}

	// rate limiter with 20% burstable rate
	var rateLimiter = make(chan time.Time, *config.RateLimit/20)
	go func() {
		for t := range time.Tick(time.Minute / time.Duration(*config.RateLimit)) {
			rateLimiter <- t
		}
	}()

	clientInstance = &Client{
		client:      http.DefaultClient,
		rateLimiter: rateLimiter,
		apiToken:    *config.ApiToken,
		baseUrl:     *config.EnvironmentUrl,
		scopes:      nil,
	}

	clientInstance.loadScopes()

	return clientInstance
}

func (at *Client) rateLimit() {
	<-at.rateLimiter
}

func (at *Client) Get(config *RequestConfig, target interface{}) error {
	at.rateLimit()

	// prepare the request URL
	req, err := at.createAuthorizedRequest(fmt.Sprintf("%s/api/v2/%s", at.baseUrl, config.Endpoint))
	if err != nil {
		return err
	}
	at.setQueryParams(req, config)

	// do the call
	err = at.do(req, target)
	if err != nil {
		return err
	}

	return nil
}

func (at *Client) createAuthorizedRequest(apiUrl string) (*http.Request, error) {
	log.Println(fmt.Sprintf("Resource URL: %s", apiUrl))

	// make a new request
	req, err := http.NewRequestWithContext(context.Background(), "GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	// set headers and query params
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", at.apiToken))

	return req, nil
}

func (at *Client) setQueryParams(req *http.Request, config *RequestConfig) {
	// set pagination params
	if config.Pagination != nil {
		config.Params.Set("pg[offset]", strconv.Itoa(config.Pagination.Offset))
		config.Params.Set("pg[limit]", strconv.Itoa(config.Pagination.Limit))
	}

	// encode params
	req.URL.RawQuery = config.Params.Encode()

	log.Println(fmt.Sprintf("Query Params: %s", req.URL.RawQuery))
}

func (at *Client) do(req *http.Request, response interface{}) error {
	var reqUrl = req.URL.RequestURI()

	// do the call
	resp, err := at.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failure on %s: %w", reqUrl, err)
	}
	defer resp.Body.Close()

	// handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return NewHttpError(reqUrl, resp)
	}

	// read response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read error on response for %s: %w", reqUrl, err)
	}

	// parse the body
	err = json.Unmarshal(b, response)
	if err != nil {
		return fmt.Errorf("JSON decode failed on %s: %s error: %w", reqUrl, b, err)
	}

	return nil
}

func (at *Client) loadScopes() {
	var config = NewRequestConfig("users/me/api-tokens")
	var result = &ApiTokenListResponse{}
	err := at.Get(config, result)
	if err == nil {
		for _, token := range result.ApiTokens {
			if at.IsTokenActive(token.Token) {
				at.scopes = &token.Scope
			}
		}
	}
}

func (at *Client) IsTokenActive(maskedToken string) bool {
	var parts = strings.Split(maskedToken, "-")
	return len(parts) > 0 && strings.HasPrefix(at.apiToken, parts[0])
}

func (at *Client) scopesLoaded() bool {
	return at.scopes != nil
}

func (at *Client) hasScope(scope string) bool {
	if at.scopes == nil {
		return false
	}

	for _, v := range *at.scopes {
		if v == scope {
			return true
		}
	}

	return false
}

func (at *Client) handleKnownErrors(err error, scope string) error {
	var httpErr = getHttpError(err)
	if httpErr == nil {
		return err
	}

	// 403 Forbidden or 404 Not Found
	if httpErr.StatusCode == 403 || httpErr.StatusCode == 404 {
		return fmt.Errorf(`We couldn't fetch the resource you requested. You either don't have access to it, or it doesn't exist.`)
	}

	// 401 Unauthorized
	if httpErr.StatusCode == 401 {
		if at.scopesLoaded() && !at.hasScope(scope) {
			return fmt.Errorf(`We couldn't fetch the resource you requested, because your API Token is missing "%s" in the enabled scopes - create a new API Token and add this scope to the list, please.`, scope)
		} else {
			return fmt.Errorf(`We couldn't fetch the resource you requested. This might be caused by "%s" scope not being enabled. Check your API Token settings in Make, please.`, scope)
		}
	}

	return httpErr
}
