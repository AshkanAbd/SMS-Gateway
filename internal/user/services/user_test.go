package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/mocks"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/services"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	inputUser := models.User{
		Name: "AshkanAbd",
	}

	t.Run("should create user", func(t *testing.T) {
		expectedUser := models.User{
			Entity: &shared.Entity{
				ID: "1",
			},
			Name:    "AshkanAbd",
			Balance: 0,
			CreateDate: &shared.CreateDate{
				CreatedAt: time.Now(),
			},
			UpdateDate: &shared.UpdateDate{
				UpdatedAt: time.Now(),
			},
		}

		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			CreateUser(ctx, inputUser).
			Return(expectedUser, nil).
			Once()

		service := services.NewUserService(mockRepo)
		actualUser, actualErr := service.CreateUser(ctx, inputUser)

		assert.NoError(t, actualErr)
		assert.Equal(t, expectedUser, actualUser)
	})

	t.Run("should return error if name is empty", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			CreateUser(ctx, inputUser).
			Return(models.User{}, models.EmptyNameError).
			Once()

		service := services.NewUserService(mockRepo)
		actualUser, actualErr := service.CreateUser(ctx, inputUser)

		assert.Error(t, actualErr)
		assert.Equal(t, models.EmptyNameError, actualErr)
		assert.Equal(t, models.User{}, actualUser)
	})
}

func TestUserService_GetUser(t *testing.T) {
	inputID := "1"

	t.Run("should return user when exist", func(t *testing.T) {
		expectedUser := models.User{
			Entity: &shared.Entity{
				ID: "1",
			},
			Name:    "AshkanAbd",
			Balance: 0,
			CreateDate: &shared.CreateDate{
				CreatedAt: time.Now(),
			},
			UpdateDate: &shared.UpdateDate{
				UpdatedAt: time.Now(),
			},
		}

		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			GetUser(ctx, inputID).
			Return(expectedUser, nil).
			Once()

		service := services.NewUserService(mockRepo)
		actualUser, actualErr := service.GetUser(ctx, inputID)

		assert.NoError(t, actualErr)
		assert.Equal(t, expectedUser, actualUser)
	})

	t.Run("should return error if user not exists", func(t *testing.T) {
		repoErr := models.UserNotExistError

		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			GetUser(ctx, inputID).
			Return(models.User{}, repoErr).
			Once()

		service := services.NewUserService(mockRepo)
		actualUser, actualErr := service.GetUser(ctx, inputID)

		assert.Error(t, actualErr)
		assert.Equal(t, repoErr, actualErr)
		assert.Equal(t, models.User{}, actualUser)
	})
}

func TestUserService_IncreaseUserBalance(t *testing.T) {
	inputID := "1"

	t.Run("should update user", func(t *testing.T) {
		inputAmount := int64(100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			UpdateUserBalance(ctx, inputID, inputAmount).
			Return(nil).
			Once()

		service := services.NewUserService(mockRepo)
		actualErr := service.IncreaseUserBalance(ctx, inputID, inputAmount)

		assert.NoError(t, actualErr)
	})

	t.Run("should do nothing when amount eq 0", func(t *testing.T) {
		inputAmount := int64(0)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		service := services.NewUserService(mockRepo)
		actualErr := service.IncreaseUserBalance(ctx, inputID, inputAmount)

		assert.NoError(t, actualErr)
	})

	t.Run("should return error when user not exists", func(t *testing.T) {
		inputAmount := int64(100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			UpdateUserBalance(ctx, inputID, inputAmount).
			Return(models.UserNotExistError).
			Once()

		service := services.NewUserService(mockRepo)
		actualErr := service.IncreaseUserBalance(ctx, inputID, inputAmount)

		assert.Error(t, actualErr)
		assert.Equal(t, models.UserNotExistError, actualErr)
	})

	t.Run("should return error when user exists but amount lt 0", func(t *testing.T) {
		inputAmount := int64(-100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		service := services.NewUserService(mockRepo)
		actualErr := service.IncreaseUserBalance(ctx, inputID, inputAmount)

		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidBalanceError, actualErr)
	})
}

func TestUserService_DecreaseUserBalance(t *testing.T) {
	inputID := "1"

	t.Run("should update user", func(t *testing.T) {
		inputAmount := int64(100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			UpdateUserBalance(ctx, inputID, inputAmount).
			Return(nil).
			Once()

		service := services.NewUserService(mockRepo)
		actualErr := service.DecreaseUserBalance(ctx, inputID, inputAmount)

		assert.NoError(t, actualErr)
	})

	t.Run("should do nothing when amount eq 0", func(t *testing.T) {
		inputAmount := int64(0)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		service := services.NewUserService(mockRepo)
		actualErr := service.DecreaseUserBalance(ctx, inputID, inputAmount)

		assert.NoError(t, actualErr)
	})

	t.Run("should return error when user not exists", func(t *testing.T) {
		inputAmount := int64(100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			UpdateUserBalance(ctx, inputID, inputAmount).
			Return(models.UserNotExistError).
			Once()

		service := services.NewUserService(mockRepo)
		actualErr := service.DecreaseUserBalance(ctx, inputID, inputAmount)

		assert.Error(t, actualErr)
		assert.Equal(t, models.UserNotExistError, actualErr)
	})

	t.Run("should return error when user exists but amount lt 0", func(t *testing.T) {
		inputAmount := int64(-100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		service := services.NewUserService(mockRepo)
		actualErr := service.DecreaseUserBalance(ctx, inputID, inputAmount)

		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidBalanceError, actualErr)
	})

	t.Run("should return error when user exists but final balance lt 0", func(t *testing.T) {
		inputAmount := int64(100)
		ctx := context.Background()
		mockRepo := mocks.NewMockIUserRepository(t)

		mockRepo.EXPECT().
			UpdateUserBalance(ctx, inputID, inputAmount).
			Return(models.InsufficientBalanceError).
			Once()

		service := services.NewUserService(mockRepo)
		actualErr := service.DecreaseUserBalance(ctx, inputID, inputAmount)

		assert.Error(t, actualErr)
		assert.Equal(t, models.InsufficientBalanceError, actualErr)
	})
}
