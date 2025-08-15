package pgsql

import (
	"fmt"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/common"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
)

type smsEntity struct {
	ID        uint
	UserId    uint
	Content   string
	Receiver  string
	Cost      int
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *smsEntity) TableName() string {
	return "messages"
}

func fromMessage(s models.Sms) smsEntity {
	se := smsEntity{
		UserId:   common.ParseUIntWithFallback(s.UserId, 0),
		Content:  s.Content,
		Receiver: s.Receiver,
		Cost:     s.Cost,
		Status:   int(s.Status),
	}

	if s.Entity != nil {
		se.ID = common.ParseUIntWithFallback(s.ID, 0)
	}
	if s.CreateDate != nil {
		se.CreatedAt = s.CreatedAt
	}
	if s.UpdateDate != nil {
		se.UpdatedAt = s.UpdatedAt
	}

	return se
}

func toMessage(se smsEntity) models.Sms {
	return models.Sms{
		Entity: &shared.Entity{
			ID: fmt.Sprintf("%d", se.ID),
		},
		CreateDate: &shared.CreateDate{
			CreatedAt: se.CreatedAt,
		},
		UpdateDate: &shared.UpdateDate{
			UpdatedAt: se.UpdatedAt,
		},
		UserId:   fmt.Sprintf("%d", se.UserId),
		Content:  se.Content,
		Receiver: se.Receiver,
		Cost:     se.Cost,
		Status:   models.SmsStatus(se.Status),
	}
}
