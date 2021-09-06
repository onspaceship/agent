package handler

import (
	"encoding/json"

	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
)

type deliveryPayload struct {
	AppId      string `json:"app_id"`
	DeliveryId string `json:"delivery_id"`
	TeamHandle string `json:"image_uri"`
}

func handleDelivery(jsonPayload []byte, options *config.SocketOptions) {
	var payload deliveryPayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	log.WithField("payload", payload).Info("Handling delivery")
}

func init() {
	handlers["delivery"] = handleDelivery
}
