package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ue := fromUser(user)

	res := r.conn.WithContext(ctx).Create(&ue)
	if res.Error != nil {
		return models.User{}, res.Error
	}

	user.Entity = &shared.Entity{
		ID: fmt.Sprintf("%d", ue.ID),
	}
	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, id string) (models.User, error) {
	ue := UserEntity{}

	err := r.conn.WithContext(ctx).First(&ue, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.UserNotExistError
		}

		return models.User{}, err
	}

	return toUser(ue), nil
}

func (r *Repository) DecreaseUserBalance(ctx context.Context, id string, amount int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) IncreaseUserBalance(ctx context.Context, id string, amount int64) error {
	//TODO implement me
	panic("implement me")
}
