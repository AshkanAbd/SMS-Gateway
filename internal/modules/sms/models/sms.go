package models

import "github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"

type SmsStatus int

const (
	StatusScheduled SmsStatus = iota
	StatusEnqueued
	StatusSent
	StatusFailed
)

type Sms struct {
	*shared.Entity
	*shared.CreateDate
	*shared.UpdateDate

	UserId   string
	Content  string
	Receiver string
	Cost     int
	Status   SmsStatus
}
