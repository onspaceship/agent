package handler

import (
	"encoding/json"

	"github.com/onspaceship/agent/pkg/config"

	"github.com/apex/log"
)

var handlers = map[string]func(payload []byte, options *config.SocketOptions){}

func Handle(event string, payload interface{}, options *config.SocketOptions) {
	if handler, ok := handlers[event]; ok {

		jsonPayload, _ := json.Marshal(payload)

		go handler(jsonPayload, options)
	} else {
		log.WithField("event", event).WithField("payload", payload).Debug("Unhandled message")
	}
}
