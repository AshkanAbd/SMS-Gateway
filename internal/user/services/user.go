package services

import "github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/repositories"

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
