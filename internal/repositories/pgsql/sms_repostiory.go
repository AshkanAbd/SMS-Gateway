package pgsql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repository) GetMessagesByUserId(ctx context.Context, userId string, skip int, limit int, desc bool) ([]models.Sms, error) {
	var ses []smsEntity

	query := r.conn.WithContext(ctx).
		Where("user_id = ?", userId).
		Limit(limit).
		Offset(skip)

	if desc {
		query = query.Order("created_at DESC, id DESC")
	} else {
		query = query.Order("created_at ASC, id ASC")
	}

	err := query.Find(&ses).Error
	if err != nil {
		return nil, err
	}

	ss := make([]models.Sms, len(ses))
	for i := range ses {
		ss[i] = toMessage(ses[i])
	}

	return ss, nil
}

func (r *Repository) SetMessageAsFailed(ctx context.Context, id string) (models.Sms, error) {
	se := smsEntity{}

	res := r.conn.WithContext(ctx).
		Model(&se).
		Clauses(clause.Returning{}).
		Where("id = ? AND status = ?", id, models.StatusEnqueued).
		Updates(map[string]any{
			"status":     models.StatusFailed,
			"updated_at": time.Now(),
		})

	if res.Error != nil {
		return models.Sms{}, res.Error
	}

	if res.RowsAffected == 0 {
		return models.Sms{}, models.MessageNotExistError
	}

	return toMessage(se), nil
}

func (r *Repository) SetMessageAsSent(ctx context.Context, id string) (models.Sms, error) {
	se := smsEntity{}

	res := r.conn.WithContext(ctx).
		Model(&se).
		Clauses(clause.Returning{}).
		Where("id = ? AND status = ?", id, models.StatusEnqueued).
		Updates(map[string]any{
			"status":     models.StatusSent,
			"updated_at": time.Now(),
		})

	if res.Error != nil {
		return models.Sms{}, res.Error
	}

	if res.RowsAffected == 0 {
		return models.Sms{}, models.MessageNotExistError
	}

	return toMessage(se), nil
}

func (r *Repository) EnqueueMessages(ctx context.Context, count int) ([]models.Sms, error) {
	var ses []smsEntity

	err := r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.WithContext(ctx).
			Model(&ses).
			Where("id IN (?)",
				tx.WithContext(ctx).
					Model(&smsEntity{}).
					Select("id").
					Where("status", models.StatusScheduled).
					Order("created_at ASC, id ASC").
					Limit(count).
					Clauses(clause.Locking{
						Strength: "UPDATE",
						Options:  "SKIP LOCKED",
					}),
			).Clauses(clause.Returning{}).
			Updates(map[string]any{
				"status":     models.StatusEnqueued,
				"updated_at": time.Now(),
			})

		if res.Error != nil {
			return res.Error
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return nil, err
	}

	ss := make([]models.Sms, len(ses))
	for i := range ses {
		ss[i] = toMessage(ses[i])
	}

	return ss, nil
}

func (r *Repository) RescheduledMessages(ctx context.Context, ids []string) error {
	err := r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.WithContext(ctx).
			Model(&smsEntity{}).
			Clauses(clause.Returning{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"status":     models.StatusScheduled,
				"updated_at": time.Now(),
			})

		if res.Error != nil {
			return res.Error
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return err
	}

	return nil
}
