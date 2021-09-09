package socket

import (
	"encoding/json"

	"github.com/apex/log"
)

const (
	AppIdAnnotation      = "onspaceship.com/app-id"
	DeliveryIdAnnotation = "onspaceship.com/delivery-id"
)

type deliveryPayload struct {
	AppId      string `json:"app_id"`
	DeliveryId string `json:"delivery_id"`
	ImageURI   string `json:"image_uri"`
}

func handleDelivery(jsonPayload []byte, socket *socket) {
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
