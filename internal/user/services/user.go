package services

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/repositories"
)

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
