package models

import "errors"

var (
	EmptyContentError  = errors.New("content is empty")
	EmptyReceiverError = errors.New("receiver is empty")
)
