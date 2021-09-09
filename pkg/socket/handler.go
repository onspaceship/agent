package socket

import (
	"encoding/json"

	"github.com/apex/log"
)

var handlers = map[string]func(payload []byte, socket *socket){}

func (socket *socket) handleEvent(event string, payload interface{}) {
	if handler, ok := handlers[event]; ok {

		jsonPayload, _ := json.Marshal(payload)

		go handler(jsonPayload, socket)
	} else {
		log.WithField("event", event).WithField("payload", payload).Debug("Unhandled message")
	}
}
