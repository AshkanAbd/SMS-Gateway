package models

import "github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"

type User struct {
	*shared.Entity
	*shared.CreateDate
	*shared.UpdateDate

	Name    string
	Balance int64
}
