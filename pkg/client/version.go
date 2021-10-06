package client

type AgentVersion struct {
	Version string `json:"version"`
}

func (client *Client) GetVersion() (AgentVersion, error) {
	var version AgentVersion

	err := client.Get(client.corePath("/agent/version"), &version)

	return version, err
}
