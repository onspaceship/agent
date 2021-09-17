package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
)

type Options = config.ClientOptions

type Client struct {
	*Options
}

func NewClient() *Client {
	options, err := config.NewClientOptions()
	if err != nil {
		log.WithError(err).Fatal("failed to configure API client")
	}

	return &Client{Options: options}
}

func (client *Client) Put(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.AgentId))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (client *Client) corePath(path string, tokens ...interface{}) string {
	url, _ := client.CoreBaseURL.Parse(fmt.Sprintf(path, tokens...))
	return url.String()
}
