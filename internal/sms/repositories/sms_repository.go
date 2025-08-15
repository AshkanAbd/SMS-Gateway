package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
)

type ISmsRepository interface {
	CreateScheduleMessages(ctx context.Context, msgs []models.Sms) error
	GetMessagesByUserId(ctx context.Context, userId string) ([]models.Sms, error)
	EnqueueMessages(ctx context.Context, count int) ([]models.Sms, error)
	RescheduledMessages(ctx context.Context, ids []string) error
	SetMessageAsFailed(ctx context.Context, id string) (models.Sms, error)
	SetMessageAsSent(ctx context.Context, id string) (models.Sms, error)
}
