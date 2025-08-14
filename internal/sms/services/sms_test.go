package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/mocks"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/services"
	"github.com/stretchr/testify/assert"
)

func TestSmsService_ScheduleMessage(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should schedule sms", func(t *testing.T) {
		ctx := context.Background()
		inputMsgs := []models.Sms{
			{
				Content:  "Test",
				Receiver: "09123456789",
				Cost:     100,
			},
		}
		inputUserId := "1"

		sendingMsgs := make([]models.Sms, len(inputMsgs))
		for i, msg := range inputMsgs {
			msg.UserId = inputUserId
			msg.Status = models.StatusScheduled
			sendingMsgs[i] = msg
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			CreateScheduleMessages(ctx, sendingMsgs).
			Return(nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualErr := service.ScheduleSms(ctx, inputUserId, inputMsgs)
		assert.NoError(t, actualErr)
	})

	t.Run("should return error when can not schedule", func(t *testing.T) {
		ctx := context.Background()
		inputMsgs := []models.Sms{
			{
				Content:  "Test",
				Receiver: "09123456789",
				Cost:     100,
			},
		}
		inputUserId := "1"

		sendingMsgs := make([]models.Sms, len(inputMsgs))
		for i, msg := range inputMsgs {
			msg.UserId = inputUserId
			msg.Status = models.StatusScheduled
			sendingMsgs[i] = msg
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			CreateScheduleMessages(ctx, sendingMsgs).
			Return(fmt.Errorf("error")).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualErr := service.ScheduleSms(ctx, inputUserId, inputMsgs)
		assert.Error(t, actualErr)
	})

	t.Run("should return EmptyContentError when can not schedule", func(t *testing.T) {
		ctx := context.Background()
		inputMsgs := []models.Sms{
			{
				Content:  "",
				Receiver: "09123456789",
				Cost:     100,
			},
		}
		inputUserId := "1"

		sendingMsgs := make([]models.Sms, len(inputMsgs))
		for i, msg := range inputMsgs {
			msg.UserId = inputUserId
			msg.Status = models.StatusScheduled
			sendingMsgs[i] = msg
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			CreateScheduleMessages(ctx, sendingMsgs).
			Return(models.EmptyReceiverError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualErr := service.ScheduleSms(ctx, inputUserId, inputMsgs)
		assert.Error(t, actualErr)
		assert.Equal(t, models.EmptyReceiverError, actualErr)
	})

	t.Run("should return EmptyReceiverError when can not schedule", func(t *testing.T) {
		ctx := context.Background()
		inputMsgs := []models.Sms{
			{
				Content:  "Test",
				Receiver: "",
				Cost:     100,
			},
		}
		inputUserId := "1"

		sendingMsgs := make([]models.Sms, len(inputMsgs))
		for i, msg := range inputMsgs {
			msg.UserId = inputUserId
			msg.Status = models.StatusScheduled
			sendingMsgs[i] = msg
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			CreateScheduleMessages(ctx, sendingMsgs).
			Return(models.EmptyReceiverError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualErr := service.ScheduleSms(ctx, inputUserId, inputMsgs)
		assert.Error(t, actualErr)
		assert.Equal(t, models.EmptyReceiverError, actualErr)
	})
}

func TestSmsService_GetUserMessages(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should return user messages", func(t *testing.T) {
		ctx := context.Background()
		inputUserId := "1"

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		expectedMsgs := []models.Sms{
			{
				Entity: &shared.Entity{
					ID: "1",
				},
				CreateDate: &shared.CreateDate{
					CreatedAt: time.Now(),
				},
				UpdateDate: &shared.UpdateDate{
					UpdatedAt: time.Now(),
				},
				UserId:   "1",
				Content:  "TestContent",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		mockRepo.EXPECT().
			GetMessagesByUserId(ctx, inputUserId).
			Return(expectedMsgs, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsgs, actualErr := service.GetUserSms(ctx, inputUserId)
		assert.NoError(t, actualErr)
		assert.Equal(t, len(expectedMsgs), len(actualMsgs))
		for i := 0; i < len(expectedMsgs); i++ {
			assert.NotNil(t, actualMsgs[i].Entity)
			assert.NotNil(t, expectedMsgs[i].Entity)
			assert.NotNil(t, actualMsgs[i].CreateDate)
			assert.NotNil(t, expectedMsgs[i].CreateDate)
			assert.NotNil(t, actualMsgs[i].UpdatedAt)
			assert.NotNil(t, expectedMsgs[i].UpdatedAt)
			assert.Equal(t, expectedMsgs[i].ID, actualMsgs[i].ID)
			assert.Equal(t, expectedMsgs[i].CreatedAt, actualMsgs[i].CreatedAt)
			assert.Equal(t, expectedMsgs[i].UpdatedAt, actualMsgs[i].UpdatedAt)
			assert.Equal(t, expectedMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, expectedMsgs[i].UserId, actualMsgs[i].UserId)
			assert.Equal(t, expectedMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, expectedMsgs[i].Receiver, actualMsgs[i].Receiver)
			assert.Equal(t, expectedMsgs[i].Cost, actualMsgs[i].Cost)
			assert.Equal(t, expectedMsgs[i].Status, actualMsgs[i].Status)
		}
	})

	t.Run("should return error when can not get messages", func(t *testing.T) {
		ctx := context.Background()
		inputUserId := "1"

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			GetMessagesByUserId(ctx, inputUserId).
			Return(nil, fmt.Errorf("error")).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsgs, actualErr := service.GetUserSms(ctx, inputUserId)
		assert.Error(t, actualErr)
		assert.Nil(t, actualMsgs)
	})
}

func TestSmsService_EnqueueEarliest(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should return enqueue count", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		expectedCount := 4

		mockRepo.EXPECT().
			EnqueueEarliestMessage(ctx).
			Return(1, nil).
			Times(expectedCount)

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, 4)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedCount, actualCount)
	})

	t.Run("should return lt count when not enough message available", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		expectedCount := 3

		mockRepo.EXPECT().
			EnqueueEarliestMessage(ctx).
			Return(1, nil).
			Times(3)
		mockRepo.EXPECT().
			EnqueueEarliestMessage(ctx).
			Return(0, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, 4)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedCount, actualCount)
	})

	t.Run("should return error and stop executing when an error occurred", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			EnqueueEarliestMessage(ctx).
			Return(1, nil).
			Times(2)
		mockRepo.EXPECT().
			EnqueueEarliestMessage(ctx).
			Return(0, fmt.Errorf("error")).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, 10)
		assert.Error(t, actualErr)
		assert.Equal(t, "error", actualErr.Error())
		assert.Equal(t, 0, actualCount)
	})
}
