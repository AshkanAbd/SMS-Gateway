package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
)

type ISmsQueue interface {
	Enqueue(ctx context.Context, msg models.Sms) error
	GetLength(ctx context.Context) (int, error)
	Pop(ctx context.Context) (models.Sms, error)
}
