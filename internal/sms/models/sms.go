package models

type SmsStatus int

const (
	Enqueued SmsStatus = iota
	Sending
	Sent
	Failed
)

type Message struct {
	ID       string
	UserId   string
	Content  string
	Receiver string
	Cost     int
	Status   SmsStatus
}
