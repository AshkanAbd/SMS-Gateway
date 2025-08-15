package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/mocks"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
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
	cfg := services.SmsServiceConfig{
		QueueCapacity: 100,
	}

	t.Run("should return enqueue count", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		msgs := []models.Sms{
			{
				Entity:     &shared.Entity{ID: "1"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 1",
				Receiver:   "09123456789",
				Cost:       100,
				Status:     models.StatusEnqueued,
			},
			{
				Entity:     &shared.Entity{ID: "2"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "2",
				Content:    "Test Content 2",
				Receiver:   "09123456788",
				Cost:       200,
				Status:     models.StatusEnqueued,
			}, {
				Entity:     &shared.Entity{ID: "3"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "3",
				Content:    "Test Content 3",
				Receiver:   "09123456787",
				Cost:       100,
				Status:     models.StatusEnqueued,
			},
			{
				Entity:     &shared.Entity{ID: "4"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "4",
				Content:    "Test Content 4",
				Receiver:   "09123456786",
				Cost:       200,
				Status:     models.StatusEnqueued,
			},
		}

		mockRepo.EXPECT().
			EnqueueMessages(ctx, len(msgs)).
			Return(msgs, nil).
			Once()

		mockQueue.EXPECT().
			Enqueue(ctx, msgs).
			Return(nil).
			Once()

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(10, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, len(msgs))
		assert.NoError(t, actualErr)
		assert.Equal(t, len(msgs), actualCount)
	})

	t.Run("should return lt count when not enough message available", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		msgs := []models.Sms{
			{
				Entity:     &shared.Entity{ID: "1"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 1",
				Receiver:   "09123456789",
				Cost:       100,
				Status:     models.StatusEnqueued,
			},
			{
				Entity:     &shared.Entity{ID: "2"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "2",
				Content:    "Test Content 2",
				Receiver:   "09123456788",
				Cost:       200,
				Status:     models.StatusEnqueued,
			},
		}

		countInput := 4
		expectedCount := len(msgs)

		mockRepo.EXPECT().
			EnqueueMessages(ctx, countInput).
			Return(msgs, nil).
			Once()

		mockQueue.EXPECT().
			Enqueue(ctx, msgs).
			Return(nil).
			Once()

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(10, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, countInput)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedCount, actualCount)
	})

	t.Run("should return db error when can not enqueue message in db", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		countInput := 4

		mockRepo.EXPECT().
			EnqueueMessages(ctx, countInput).
			Return(nil, fmt.Errorf("db error")).
			Once()

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(10, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, countInput)
		assert.Error(t, actualErr)
		assert.Equal(t, "db error", actualErr.Error())
		assert.Equal(t, 0, actualCount)
	})

	t.Run("should return queue error and reschedule when can not enqueue message to queue", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		msgs := []models.Sms{
			{
				Entity:     &shared.Entity{ID: "1"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 1",
				Receiver:   "09123456789",
				Cost:       100,
				Status:     models.StatusEnqueued,
			},
			{
				Entity:     &shared.Entity{ID: "2"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "2",
				Content:    "Test Content 2",
				Receiver:   "09123456788",
				Cost:       200,
				Status:     models.StatusEnqueued,
			},
		}
		ids := make([]string, len(msgs))
		for i := range msgs {
			ids[i] = msgs[i].ID
		}

		countInput := 4

		mockRepo.EXPECT().
			EnqueueMessages(ctx, countInput).
			Return(msgs, nil).
			Once()

		mockQueue.EXPECT().
			Enqueue(ctx, msgs).
			Return(fmt.Errorf("queue error")).
			Once()

		mockRepo.EXPECT().
			RescheduledMessages(ctx, ids).
			Return(nil).
			Once()

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(10, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, countInput)
		assert.Error(t, actualErr)
		assert.Equal(t, "queue error", actualErr.Error())
		assert.Equal(t, 0, actualCount)
	})

	t.Run("should return db error when can not reschedule", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		msgs := []models.Sms{
			{
				Entity:     &shared.Entity{ID: "1"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 1",
				Receiver:   "09123456789",
				Cost:       100,
				Status:     models.StatusEnqueued,
			},
			{
				Entity:     &shared.Entity{ID: "2"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "2",
				Content:    "Test Content 2",
				Receiver:   "09123456788",
				Cost:       200,
				Status:     models.StatusEnqueued,
			},
		}
		ids := make([]string, len(msgs))
		for i := range msgs {
			ids[i] = msgs[i].ID
		}

		countInput := 4

		mockRepo.EXPECT().
			EnqueueMessages(ctx, countInput).
			Return(msgs, nil).
			Once()

		mockQueue.EXPECT().
			Enqueue(ctx, msgs).
			Return(fmt.Errorf("queue error")).
			Once()

		mockRepo.EXPECT().
			RescheduledMessages(ctx, ids).
			Return(fmt.Errorf("db error")).
			Once()

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(10, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, countInput)
		assert.Error(t, actualErr)
		assert.Equal(t, "db error", actualErr.Error())
		assert.Equal(t, 0, actualCount)
	})

	t.Run("should return InvalidQueueError when queue is not valid", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(0, models.InvalidQueueError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, 1)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidQueueError, actualErr)
		assert.Equal(t, 0, actualCount)
	})

	t.Run("should return NoCapacityInQueueError when queue is full", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			GetLength(ctx).
			Return(100, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualCount, actualErr := service.EnqueueEarliest(ctx, 1)
		assert.Error(t, actualErr)
		assert.Equal(t, models.NoCapacityInQueueError, actualErr)
		assert.Equal(t, 0, actualCount)
	})
}

func TestSmsService_SetMessageAsFailed(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should set message as failed", func(t *testing.T) {
		ctx := context.Background()
		expected := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusFailed,
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			SetMessageAsFailed(ctx, expected.ID).
			Return(expected, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SetMessageAsFailed(ctx, expected.ID)
		assert.NoError(t, actualErr)
		assert.Equal(t, expected.ID, actualMsg.ID)
		assert.Equal(t, expected.CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, expected.UserId, actualMsg.UserId)
		assert.Equal(t, expected.Content, actualMsg.Content)
		assert.Equal(t, expected.Receiver, actualMsg.Receiver)
		assert.Equal(t, expected.Cost, actualMsg.Cost)
		assert.Equal(t, expected.Status, actualMsg.Status)
	})

	t.Run("should return MessageNotExistError when message not exists", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			SetMessageAsFailed(ctx, "1").
			Return(models.Sms{}, models.MessageNotExistError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SetMessageAsFailed(ctx, "1")
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
		assert.Equal(t, models.Sms{}, actualMsg)
	})
}

func TestSmsService_SetMessageAsSent(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should set message as sent", func(t *testing.T) {
		ctx := context.Background()
		expectedMsg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusSent,
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			SetMessageAsSent(ctx, expectedMsg.ID).
			Return(expectedMsg, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SetMessageAsSent(ctx, expectedMsg.ID)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedMsg.ID, actualMsg.ID)
		assert.Equal(t, expectedMsg.CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})

	t.Run("should return MessageNotExistError when message not exists", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockRepo.EXPECT().
			SetMessageAsSent(ctx, "1").
			Return(models.Sms{}, models.MessageNotExistError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SetMessageAsSent(ctx, "1")
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
		assert.Equal(t, models.Sms{}, actualMsg)
	})
}

func TestSmsService_SendFromQueue(t *testing.T) {
	cfg := services.SmsServiceConfig{}

	t.Run("should read message from queue and send it and update its status to sent", func(t *testing.T) {
		ctx := context.Background()
		msg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}
		expectedMsg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusSent,
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			Pop(ctx).
			Return(msg, nil).
			Once()

		mockSender.EXPECT().
			Send(ctx, msg).
			Return(nil).
			Once()

		mockRepo.EXPECT().
			SetMessageAsSent(ctx, msg.ID).
			Return(expectedMsg, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SendFromQueue(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedMsg.ID, actualMsg.ID)
		assert.Equal(t, expectedMsg.CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})

	t.Run("should read message from queue and can not send it and update its status to failed", func(t *testing.T) {
		ctx := context.Background()
		msg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}
		expectedMsg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusFailed,
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			Pop(ctx).
			Return(msg, nil).
			Once()

		mockSender.EXPECT().
			Send(ctx, msg).
			Return(models.SendError).
			Once()

		mockRepo.EXPECT().
			SetMessageAsFailed(ctx, msg.ID).
			Return(expectedMsg, nil).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SendFromQueue(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedMsg.ID, actualMsg.ID)
		assert.Equal(t, expectedMsg.CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})

	t.Run("should return InvalidQueueError when can not read from queue", func(t *testing.T) {
		ctx := context.Background()

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			Pop(ctx).
			Return(models.Sms{}, models.InvalidQueueError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SendFromQueue(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidQueueError, actualErr)
		assert.Equal(t, models.Sms{}, actualMsg)
	})

	t.Run("should return MessageNotExistError when message not exist", func(t *testing.T) {
		ctx := context.Background()
		msg := models.Sms{
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
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}

		mockQueue := mocks.NewMockISmsQueue(t)
		mockSender := mocks.NewMockISmsSender(t)
		mockRepo := mocks.NewMockISmsRepository(t)

		mockQueue.EXPECT().
			Pop(ctx).
			Return(msg, nil).
			Once()

		mockSender.EXPECT().
			Send(ctx, msg).
			Return(nil).
			Once()

		mockRepo.EXPECT().
			SetMessageAsSent(ctx, msg.ID).
			Return(models.Sms{}, models.MessageNotExistError).
			Once()

		service := services.NewSmsService(cfg, mockRepo, mockSender, mockQueue)

		actualMsg, actualErr := service.SendFromQueue(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
		assert.Equal(t, models.Sms{}, actualMsg)
	})
}
