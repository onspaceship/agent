package socket

import (
	"encoding/json"

	"github.com/onspaceship/agent/pkg/delivery"

	"github.com/apex/log"
)

type deliveryPayload struct {
	DeliveryId string `json:"delivery_id"`
	AppId      string `json:"app_id"`
	AppHandle  string `json:"app_handle"`
	TeamHandle string `json:"team_handle"`
	ImageURI   string `json:"image_uri"`
}

func handleDelivery(jsonPayload []byte, socket *socket) {
	var payload deliveryPayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	d := delivery.NewDelivery(payload.DeliveryId, payload.AppId, payload.AppHandle, payload.TeamHandle, payload.ImageURI, socket.client)
	d.Deliver()
}

func init() {
	handlers["delivery"] = handleDelivery
}
