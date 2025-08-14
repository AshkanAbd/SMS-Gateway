package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
)

type ISmsRepository interface {
	CreateScheduleMessages(ctx context.Context, msgs []models.Sms) error
	GetMessagesByUserId(ctx context.Context, userId string) ([]models.Sms, error)
	EnqueueEarliestMessage(ctx context.Context) (int, error)
}
