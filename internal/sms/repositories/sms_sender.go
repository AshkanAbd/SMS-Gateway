package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
)

type ISmsSender interface {
	Send(ctx context.Context, msg models.Message) error
}
