package client

type AgentSecret struct {
	DockerConfigJson string `json:"docker_config_json"`
}

func (client *Client) GetSecret() (AgentSecret, error) {
	var secret AgentSecret

	err := client.Get(client.corePath("/agent/secret"), &secret)

	return secret, err
}
