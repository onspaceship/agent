package config

import (
	"errors"
	"net/url"

	"github.com/spf13/viper"
)

const (
	DefaultCoreBaseURL = "https://local.onspaceship.dev:10001/"
)

type ClientOptions struct {
	AgentId     string
	CoreBaseURL *url.URL
}

func NewClientOptions() (*ClientOptions, error) {
	options := &ClientOptions{}
	err := options.Configure()

	return options, err
}

func (options *ClientOptions) Configure() error {
	options.AgentId = viper.GetString("agent_id")
	if options.AgentId == "" {
		return errors.New("agent ID must be provided")
	}

	coreBaseURL, err := url.Parse(viper.GetString("core_base_url"))
	if err != nil {
		return errors.New("invalid core_base_url")
	}
	options.CoreBaseURL = coreBaseURL

	return nil
}

func init() {
	viper.SetDefault("core_base_url", DefaultCoreBaseURL)
}
