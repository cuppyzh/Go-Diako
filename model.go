package diako

import (
	"log"
	"strconv"
	"time"

	"github.com/gookit/event"
)

type MessageRequest struct {
	Sender    string `json:"Sender"`
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

type DiakoMessageReceivedEventData struct {
	event.BasicEvent
	Message MessageRequest
}

func (e *DiakoMessageReceivedEventData) GetMessageData() MessageRequest {
	return e.Message
}

func (messageRequest MessageRequest) GetTime() time.Time {
	log.Println("Parse GetTime()")
	log.Println(messageRequest)

	parsedInt, err := strconv.ParseInt(messageRequest.Timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(parsedInt, 0)
}
