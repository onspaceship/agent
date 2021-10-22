package config

import (
	"errors"
	"net/url"

	"github.com/spf13/viper"
)

const (
	DefaultCoreBaseURL = "https://core.onspaceship.com/"
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

	viper.SetDefault("core_base_url", DefaultCoreBaseURL)
	coreBaseURL, err := url.Parse(viper.GetString("core_base_url"))
	if err != nil {
		return errors.New("invalid core_base_url")
	}
	options.CoreBaseURL = coreBaseURL

	return nil
}
