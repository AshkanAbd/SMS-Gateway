package dummy

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
)

type SmsSender struct {
}

func NewSmsSender() *SmsSender {
	return &SmsSender{}
}

func (d *SmsSender) Send(_ context.Context, _ models.Sms) error {
	return nil
}
