package makesdk

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"strings"
)

type Config struct {
	ApiToken       *string
	EnvironmentUrl *string
	RateLimit      *int
}

func NewConfig(apiToken *string, environmentUrl *string, rateLimit *int) (*Config, error) {
	var err error

	err = validateRateLimit(rateLimit)
	if err != nil {
		return nil, err
	}

	err = validateApiToken(apiToken)
	if err != nil {
		return nil, err
	}

	err = validateEnvironmentUrl(environmentUrl)
	if err != nil {
		return nil, err
	}
	var cleanEnvUrl = strings.TrimSuffix(*environmentUrl, "/")

	return &Config{
		ApiToken:       apiToken,
		EnvironmentUrl: &cleanEnvUrl,
		RateLimit:      rateLimit,
	}, err
}

func validateEnvironmentUrl(envUrl *string) error {
	if envUrl == nil {
		return fmt.Errorf("the environment URL is not defined")
	}

	u, err := url.ParseRequestURI(*envUrl)
	if err != nil {
		return fmt.Errorf("the environment URL does not seem to be a properly formatted URL")
	}

	if strings.ToLower(u.Scheme) != "https" {
		return fmt.Errorf("use HTTPS protocol for the environment URL")
	}

	return nil
}

func validateApiToken(apiToken *string) error {
	if apiToken == nil {
		return fmt.Errorf("the API Token is not defined; to get a token, visit the API tab in your Profile page in Make")
	}

	_, err := uuid.Parse(*apiToken)
	if err != nil {
		return fmt.Errorf("the API Token seems to have a wrong format; to get a token, visit the API tab in your Profile page in Make")
	}

	return nil
}

func validateRateLimit(rateLimit *int) error {
	if rateLimit != nil && *rateLimit <= 0 {
		return fmt.Errorf("the rate limit should be a positive number")
	}

	return nil
}
