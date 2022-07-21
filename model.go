package diako

type MessageRequest struct {
	Sender 	string	`json:"Sender"`
	Message	string	`json:"Message"`
}