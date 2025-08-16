package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ue := fromUser(user)

	res := r.conn.WithContext(ctx).Create(&ue)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "user_name_empty") {
			return models.User{}, models.EmptyNameError
		}
		return models.User{}, res.Error
	}

	user.Entity = &shared.Entity{
		ID: fmt.Sprintf("%d", ue.ID),
	}
	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, id string) (models.User, error) {
	ue := userEntity{}

	err := r.conn.WithContext(ctx).First(&ue, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.UserNotExistError
		}

		return models.User{}, err
	}

	return toUser(ue), nil
}

func (r *Repository) UpdateUserBalance(ctx context.Context, id string, amount int64) (int64, error) {
	var ue userEntity

	err := r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.WithContext(ctx).Model(&ue).
			Where("id = ?", id).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).
			Updates(
				map[string]any{
					"balance":    gorm.Expr("balance + ?", amount),
					"updated_at": gorm.Expr("now()"),
				},
			)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return models.UserNotExistError
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		if strings.Contains(err.Error(), "user_insufficient_balance") {
			return 0, models.InsufficientBalanceError
		}
		return 0, err
	}

	return ue.Balance, nil
}
