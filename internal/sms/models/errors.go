package models

import "errors"

var (
	EmptyContentError    = errors.New("content is empty")
	EmptyReceiverError   = errors.New("receiver is empty")
	InvalidQueueError    = errors.New("queue is invalid")
	SendError            = errors.New("failed to send message")
	MessageNotExistError = errors.New("message does not exist")
)
