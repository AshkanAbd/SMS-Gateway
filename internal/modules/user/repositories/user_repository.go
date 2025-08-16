package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	UpdateUserBalance(ctx context.Context, id string, amount int64) (int64, error)
}
