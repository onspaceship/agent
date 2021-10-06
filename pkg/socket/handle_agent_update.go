package socket

import (
	"encoding/json"

	"github.com/onspaceship/agent/pkg/update"

	"github.com/apex/log"
)

type agentUpdatePayload struct {
	Version string `json:"version"`
}

func handleAgentUpdate(jsonPayload []byte, socket *socket) {
	var payload agentUpdatePayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	log.WithField("version", payload.Version).Info("Handling agent update")

	err = update.ProcessVersionUpdate(payload.Version, socket.client)
	if err != nil {
		log.WithError(err).Info("Problem updating version")
		return
	}
}

func init() {
	handlers["agent_update"] = handleAgentUpdate
}
