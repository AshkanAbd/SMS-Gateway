package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
)

type ISmsRepository interface {
	Enqueue(ctx context.Context, msgs []models.Message) error
	Peek(ctx context.Context, count int) ([]models.Message, error)
	SetFailed(ctx context.Context, ids []int) error
	SetSent(ctx context.Context, ids []int) error
}
