package services

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/repositories"
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
	res, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}

func (u *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	res, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}

func (u *UserService) IncreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error) {
	if amount == 0 {
		return 0, nil
	}
	if amount < 0 {
		return 0, models.InvalidBalanceError
	}

	newBalance, err := u.userRepo.UpdateUserBalance(ctx, userId, amount)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (u *UserService) DecreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error) {
	if amount == 0 {
		return 0, nil
	}
	if amount < 0 {
		return 0, models.InvalidBalanceError
	}

	newBalance, err := u.userRepo.UpdateUserBalance(ctx, userId, amount)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}
