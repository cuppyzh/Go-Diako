package diako

import "github.com/gookit/event"

type MessageRequest struct {
	Sender  string `json:"Sender"`
	Message string `json:"Message"`
}

type DiakoMessageReceivedEventData struct {
	event.BasicEvent
	Message MessageRequest
}

func (e *DiakoMessageReceivedEventData) GetMessageData() MessageRequest {
	return e.Message
}
