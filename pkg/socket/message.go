package socket

type Message struct {
	Topic   string      `json:"topic"`
	Event   Event       `json:"event"`
	Payload interface{} `json:"payload"`
	Ref     int64       `json:"ref"`
}

type Event string

const (
	CloseEvent Event = "phx_close"
	ErrorEvent Event = "phx_error"
	JoinEvent  Event = "phx_join"
	ReplyEvent Event = "phx_reply"
	LeaveEvent Event = "phx_leave"
)

func isPhoenixEvent(event Event) bool {
	return event == CloseEvent || event == ErrorEvent || event == JoinEvent || event == ReplyEvent || event == LeaveEvent
}

func (socket *socket) joinTopic(topic string) {
	msg := &Message{
		Topic: topic,
		Event: JoinEvent,
		Ref:   socket.refCounter.nextRef(),
	}

	socket.conn.WriteJSON(msg)
}

func (socket *socket) sendMessage(event Event, payload interface{}) {
	msg := &Message{
		Topic:   "agent:" + socket.AgentId,
		Event:   event,
		Payload: payload,
		Ref:     socket.refCounter.nextRef(),
	}

	socket.conn.WriteJSON(msg)
}

func (socket *socket) sendHeartbeat() {
	msg := &Message{
		Topic:   "phoenix",
		Event:   "heartbeat",
		Payload: nil,
		Ref:     socket.refCounter.nextRef(),
	}

	socket.conn.WriteJSON(msg)
}
