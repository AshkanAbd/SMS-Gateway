package services

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/repositories"

	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

type IUserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	IncreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error)
	DecreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error)
}

type UserService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(
	userRepo repositories.IUserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	pkgLog.Debug("creating new user with name %s", user.Name)
	res, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		pkgLog.Error(err, "failed to create user with name %s", user.Name)
		return models.User{}, err
	}

	pkgLog.Debug("created new user with name %s id %s", user.Name, res.ID)
	return res, nil
}

func (u *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	pkgLog.Debug("getting user with id %s", id)
	res, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		pkgLog.Error(err, "failed to get user with id %s", id)
		return models.User{}, err
	}

	pkgLog.Debug("got user with id %s", id)
	return res, nil
}

func (u *UserService) IncreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error) {
	pkgLog.Debug("increasing user balance for user id %s amount %d", userId, amount)
	if amount == 0 {
		pkgLog.Debug("change amount is zero for user id %s, skipping", userId)
		return 0, nil
	}
	if amount < 0 {
		pkgLog.Error(models.InvalidBalanceError, "negative amount for user id %s", userId)
		return 0, models.InvalidBalanceError
	}

	newBalance, err := u.userRepo.UpdateUserBalance(ctx, userId, amount)
	if err != nil {
		pkgLog.Error(err, "failed to increase user balance for user id %s", userId)
		return 0, err
	}

	pkgLog.Debug("increased user balance for user id %s amount %d to %d", userId, amount, newBalance)
	return newBalance, nil
}

func (u *UserService) DecreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error) {
	pkgLog.Debug("decreasing user balance for user id %s amount %d", userId, amount)
	if amount == 0 {
		pkgLog.Debug("change amount is zero for user id %s, skipping", userId)
		return 0, nil
	}
	if amount < 0 {
		pkgLog.Error(models.InvalidBalanceError, "negative amount for user id %s", userId)
		return 0, models.InvalidBalanceError
	}

	newBalance, err := u.userRepo.UpdateUserBalance(ctx, userId, -amount)
	if err != nil {
		pkgLog.Error(err, "failed to decrease user balance for user id %s", userId)
		return 0, err
	}

	pkgLog.Debug("decreased user balance for user id %s amount %d to %d", userId, amount, newBalance)
	return newBalance, nil
}
