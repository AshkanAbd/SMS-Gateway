package repositories

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
)

type IUserRepository interface {
	AddUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, id string) (models.User, error)
	DecreaseUserBalance(ctx context.Context, id string, amount int64) error
	IncreaseUserBalance(ctx context.Context, id string, amount int64) error
}
