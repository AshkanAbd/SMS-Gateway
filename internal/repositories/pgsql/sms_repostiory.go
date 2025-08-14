package pgsql

import (
	"context"
	"strings"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateScheduleMessages(ctx context.Context, msgs []models.Sms) error {
	ses := make([]smsEntity, len(msgs))
	for i := range msgs {
		ses[i] = fromMessage(msgs[i])
	}

	err := r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(ses).
			Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "message_content_empty") {
			return models.EmptyContentError
		}
		if strings.Contains(err.Error(), "message_receiver_empty") {
			return models.EmptyReceiverError
		}
		return err
	}

	return nil
}

func (r *Repository) GetMessagesByUserId(ctx context.Context, userId string) ([]models.Sms, error) {
	var ses []smsEntity

	err := r.conn.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&ses).Error
	if err != nil {
		return nil, err
	}

	ss := make([]models.Sms, len(ses))
	for i := range ses {
		ss[i] = toMessage(ses[i])
	}

	return ss, nil
}

func (r *Repository) SetFailed(ctx context.Context, ids []int) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) SetSent(ctx context.Context, ids []int) error {
	//TODO implement me
	panic("implement me")
}
