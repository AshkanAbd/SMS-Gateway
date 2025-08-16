package smsgateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/smsgateway"
	"github.com/stretchr/testify/assert"

	smsmocks "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/mocks"
	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	usermocks "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/mocks"
	usermodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
)

func TestSmsGateway_CreateUser(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    0,
		SendWorkerCount: 0,
		MessageCost:     0,
	}

	t.Run("should create new user", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		expectedUser := usermodels.User{
			Name: "AshkanAbd",
		}

		mockUser.EXPECT().
			CreateUser(ctx, expectedUser).
			Return(expectedUser, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualUser, actualErr := smsGateway.CreateUser(ctx, expectedUser)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedUser.Name, actualUser.Name)
	})

	t.Run("should return error when can not create user", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		expectedUser := usermodels.User{
			Name: "AshkanAbd",
		}

		expectedErr := fmt.Errorf("some error")

		mockUser.EXPECT().
			CreateUser(ctx, expectedUser).
			Return(usermodels.User{}, expectedErr).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualUser, actualErr := smsGateway.CreateUser(ctx, expectedUser)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
		assert.Equal(t, usermodels.User{}, actualUser)
	})
}

func TestSmsGateway_GetUser(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    0,
		SendWorkerCount: 0,
		MessageCost:     0,
	}

	t.Run("should return user", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		expectedUser := usermodels.User{
			Entity: &shared.Entity{ID: userId},
			Name:   "AshkanAbd",
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(expectedUser, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualUser, actualErr := smsGateway.GetUser(ctx, userId)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedUser.Name, actualUser.Name)
		assert.Equal(t, expectedUser.ID, actualUser.ID)
	})

	t.Run("should return error when can not get user", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		expectedErr := fmt.Errorf("some error")

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(usermodels.User{}, expectedErr).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualUser, actualErr := smsGateway.GetUser(ctx, userId)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
		assert.Equal(t, usermodels.User{}, actualUser)
	})
}

func TestSmsGateway_GetUserMessages(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    0,
		SendWorkerCount: 0,
		MessageCost:     0,
	}

	t.Run("should return user messages", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		expectedMsgs := []smsmodels.Sms{
			{
				Entity:     &shared.Entity{ID: "1"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 1",
				Receiver:   "09123456789",
				Cost:       100,
				Status:     smsmodels.StatusScheduled,
			}, {
				Entity:     &shared.Entity{ID: "2"},
				CreateDate: &shared.CreateDate{CreatedAt: time.Now()},
				UpdateDate: &shared.UpdateDate{UpdatedAt: time.Now()},
				UserId:     "1",
				Content:    "Test Content 2",
				Receiver:   "09123456788",
				Cost:       100,
				Status:     smsmodels.StatusEnqueued,
			},
		}

		mockSms.EXPECT().
			GetUserSms(ctx, userId).
			Return(expectedMsgs, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualMsgs, actualErr := smsGateway.GetUserMessages(ctx, userId)
		assert.NoError(t, actualErr)
		assert.Equal(t, len(expectedMsgs), len(actualMsgs))
		for i := 0; i < len(expectedMsgs); i++ {
			assert.Equal(t, expectedMsgs[i].ID, actualMsgs[i].ID)
			assert.Equal(t, expectedMsgs[i].CreatedAt, actualMsgs[i].CreatedAt)
			assert.Equal(t, expectedMsgs[i].UpdatedAt, actualMsgs[i].UpdatedAt)
			assert.Equal(t, expectedMsgs[i].UserId, actualMsgs[i].UserId)
			assert.Equal(t, expectedMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, expectedMsgs[i].Receiver, actualMsgs[i].Receiver)
			assert.Equal(t, expectedMsgs[i].Cost, actualMsgs[i].Cost)
			assert.Equal(t, expectedMsgs[i].Status, actualMsgs[i].Status)
		}
	})

	t.Run("should return error when can not get user messages", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		expectedErr := fmt.Errorf("some error")

		mockSms.EXPECT().
			GetUserSms(ctx, userId).
			Return(nil, expectedErr).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualMsgs, actualErr := smsGateway.GetUserMessages(ctx, userId)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
		assert.Nil(t, actualMsgs)
	})
}

func TestSmsGateway_SendSingleMessage(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    0,
		SendWorkerCount: 0,
		MessageCost:     100,
	}

	t.Run("should schedule message", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost)).
			Return(0, nil).
			Once()

		mockSms.EXPECT().
			ScheduleSms(ctx, userId, []smsmodels.Sms{
				{
					Content:  msg.Content,
					Receiver: msg.Receiver,
					Cost:     cfg.MessageCost,
				},
			}).Return(nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendSingleMessage(ctx, userId, msg)
		assert.NoError(t, actualErr)
	})

	t.Run("should return InsufficientBalanceError when user balance is not enough", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 10,
		}

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendSingleMessage(ctx, userId, msg)
		assert.Error(t, actualErr)
		assert.Equal(t, usermodels.InsufficientBalanceError, actualErr)
	})

	t.Run("should return InsufficientBalanceError when user balance is not enough because of race", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost)).
			Return(0, usermodels.InsufficientBalanceError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendSingleMessage(ctx, userId, msg)
		assert.Error(t, actualErr)
		assert.Equal(t, usermodels.InsufficientBalanceError, actualErr)
	})

	t.Run("should return error when can not schedule message", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
		}

		expectedErr := fmt.Errorf("some error")

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost)).
			Return(0, nil).
			Once()

		mockSms.EXPECT().
			ScheduleSms(ctx, userId, []smsmodels.Sms{
				{
					Content:  msg.Content,
					Receiver: msg.Receiver,
					Cost:     cfg.MessageCost,
				},
			}).Return(expectedErr).
			Once()

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, userId, int64(cfg.MessageCost)).
			Return(0, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendSingleMessage(ctx, userId, msg)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestSmsGateway_SendBulkMessage(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    0,
		SendWorkerCount: 0,
		MessageCost:     100,
	}

	t.Run("should schedule message", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msgs := []smsmodels.Sms{
			{
				Content:  "Test Content 1",
				Receiver: "09123456789",
			}, {
				Content:  "Test Content 2",
				Receiver: "09123456788",
			},
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost*len(msgs))).
			Return(0, nil).
			Once()

		mockSms.EXPECT().
			ScheduleSms(ctx, userId, []smsmodels.Sms{
				{
					Content:  msgs[0].Content,
					Receiver: msgs[0].Receiver,
					Cost:     cfg.MessageCost,
				}, {
					Content:  msgs[1].Content,
					Receiver: msgs[1].Receiver,
					Cost:     cfg.MessageCost,
				},
			}).Return(nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendBulkMessage(ctx, userId, msgs)
		assert.NoError(t, actualErr)
	})

	t.Run("should return InsufficientBalanceError when user balance is not enough", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 10,
		}

		msgs := []smsmodels.Sms{
			{
				Content:  "Test Content 1",
				Receiver: "09123456789",
			}, {
				Content:  "Test Content 2",
				Receiver: "09123456788",
			},
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendBulkMessage(ctx, userId, msgs)
		assert.Error(t, actualErr)
		assert.Equal(t, usermodels.InsufficientBalanceError, actualErr)
	})

	t.Run("should return InsufficientBalanceError when user balance is not enough because of race", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msgs := []smsmodels.Sms{
			{
				Content:  "Test Content 1",
				Receiver: "09123456789",
			}, {
				Content:  "Test Content 2",
				Receiver: "09123456788",
			},
		}

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost*len(msgs))).
			Return(0, usermodels.InsufficientBalanceError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendBulkMessage(ctx, userId, msgs)
		assert.Error(t, actualErr)
		assert.Equal(t, usermodels.InsufficientBalanceError, actualErr)
	})

	t.Run("should return error when can not schedule message", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		userId := "1"

		user := usermodels.User{
			Entity:  &shared.Entity{ID: "1"},
			Name:    "AshkanAbd",
			Balance: 1000,
		}

		msgs := []smsmodels.Sms{
			{
				Content:  "Test Content 1",
				Receiver: "09123456789",
			}, {
				Content:  "Test Content 2",
				Receiver: "09123456788",
			},
		}

		expectedErr := fmt.Errorf("some error")

		mockUser.EXPECT().
			GetUser(ctx, userId).
			Return(user, nil).
			Once()

		mockUser.EXPECT().
			DecreaseUserBalance(ctx, userId, int64(cfg.MessageCost*len(msgs))).
			Return(0, nil).
			Once()

		mockSms.EXPECT().
			ScheduleSms(ctx, userId, []smsmodels.Sms{
				{
					Content:  msgs[0].Content,
					Receiver: msgs[0].Receiver,
					Cost:     cfg.MessageCost,
				}, {
					Content:  msgs[1].Content,
					Receiver: msgs[1].Receiver,
					Cost:     cfg.MessageCost,
				},
			}).Return(expectedErr).
			Once()

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, userId, int64(cfg.MessageCost*len(msgs))).
			Return(0, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendBulkMessage(ctx, userId, msgs)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
	})
}

func TestSmsGateway_EnqueueWorker(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    10,
		SendWorkerCount: 1,
		MessageCost:     100,
	}

	t.Run("should enqueue scheduled messages", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		expectedEnqueue := 10

		mockSms.EXPECT().
			EnqueueEarliest(ctx, cfg.EnqueueCount).
			Return(10, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualEnqueue, actualErr := smsGateway.EnqueueWorker(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedEnqueue, actualEnqueue)
	})

	t.Run("should return InvalidQueueError when queue is not valid", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		mockSms.EXPECT().
			EnqueueEarliest(ctx, cfg.EnqueueCount).
			Return(0, smsmodels.InvalidQueueError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualEnqueue, actualErr := smsGateway.EnqueueWorker(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, smsmodels.InvalidQueueError, actualErr)
		assert.Equal(t, 0, actualEnqueue)
	})

	t.Run("should return NoCapacityInQueueError when queue is full", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		mockSms.EXPECT().
			EnqueueEarliest(ctx, cfg.EnqueueCount).
			Return(0, smsmodels.NoCapacityInQueueError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualEnqueue, actualErr := smsGateway.EnqueueWorker(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, smsmodels.NoCapacityInQueueError, actualErr)
		assert.Equal(t, 0, actualEnqueue)
	})

	t.Run("should return nil when other errors occurred", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		mockSms.EXPECT().
			EnqueueEarliest(ctx, cfg.EnqueueCount).
			Return(0, fmt.Errorf("some error")).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualEnqueue, actualErr := smsGateway.EnqueueWorker(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, 0, actualEnqueue)
	})
}

func TestSmsGateway_SendWorker(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    10,
		SendWorkerCount: 1,
		MessageCost:     100,
	}

	t.Run("should send a queue message", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
			UserId:   "1",
			Status:   smsmodels.StatusSent,
		}

		mockSms.EXPECT().
			SendFromQueue(ctx).
			Return(msg, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendWorker(ctx)
		assert.NoError(t, actualErr)
	})

	t.Run("should increase user balance when can not send", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
			UserId:   "1",
			Cost:     200,
			Status:   smsmodels.StatusFailed,
		}

		mockSms.EXPECT().
			SendFromQueue(ctx).
			Return(msg, nil).
			Once()

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, msg.UserId, int64(msg.Cost)).
			Return(0, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendWorker(ctx)
		assert.NoError(t, actualErr)
	})

	t.Run("should return nil when can not increase user balance", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		msg := smsmodels.Sms{
			Content:  "Test Content 1",
			Receiver: "09123456789",
			UserId:   "1",
			Cost:     200,
			Status:   smsmodels.StatusFailed,
		}

		mockSms.EXPECT().
			SendFromQueue(ctx).
			Return(msg, nil).
			Once()

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, msg.UserId, int64(msg.Cost)).
			Return(0, fmt.Errorf("some error")).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendWorker(ctx)
		assert.NoError(t, actualErr)
	})

	t.Run("should return InvalidQueueError when queue is not valid", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		mockSms.EXPECT().
			SendFromQueue(ctx).
			Return(smsmodels.Sms{}, smsmodels.InvalidQueueError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendWorker(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, smsmodels.InvalidQueueError, actualErr)
	})

	t.Run("should return nil when message not exists", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		mockSms.EXPECT().
			SendFromQueue(ctx).
			Return(smsmodels.Sms{}, smsmodels.MessageNotExistError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualErr := smsGateway.SendWorker(ctx)
		assert.NoError(t, actualErr)
	})
}

func TestSmsGateway_IncreaseUserBalance(t *testing.T) {
	cfg := smsgateway.Config{
		EnqueueCount:    10,
		SendWorkerCount: 1,
		MessageCost:     100,
	}

	t.Run("should increase user balance", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		inputUserId := "1"
		inputAmount := int64(100)

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, inputUserId, inputAmount).
			Return(inputAmount, nil).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualBalance, actualErr := smsGateway.IncreaseUserBalance(ctx, inputUserId, inputAmount)
		assert.NoError(t, actualErr)
		assert.Equal(t, inputAmount, actualBalance)
	})

	t.Run("should return UserNotExistError when can not increase balance", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		inputUserId := "1"
		inputAmount := int64(100)

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, inputUserId, inputAmount).
			Return(0, usermodels.UserNotExistError).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualBalance, actualErr := smsGateway.IncreaseUserBalance(ctx, inputUserId, inputAmount)
		assert.Error(t, actualErr)
		assert.Equal(t, usermodels.UserNotExistError, actualErr)
		assert.Equal(t, int64(0), actualBalance)
	})

	t.Run("should return error when can not increase balance", func(t *testing.T) {
		ctx := context.Background()

		mockUser := usermocks.NewMockIUserService(t)
		mockSms := smsmocks.NewMockISmsService(t)

		inputUserId := "1"
		inputAmount := int64(100)

		expectedErr := fmt.Errorf("some error")

		mockUser.EXPECT().
			IncreaseUserBalance(ctx, inputUserId, inputAmount).
			Return(0, expectedErr).
			Once()

		smsGateway := smsgateway.NewSmsGateway(cfg, mockUser, mockSms)

		actualBalance, actualErr := smsGateway.IncreaseUserBalance(ctx, inputUserId, inputAmount)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr, actualErr)
		assert.Equal(t, int64(0), actualBalance)
	})
}
