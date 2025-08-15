package smsgateway

import (
	"context"
	"errors"
	"time"

	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	smssrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	usermodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
	usersrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/services"
	logPkg "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

type Config struct {
	EnqueueCount              int           `mapstructure:"enqueue_count"`
	FullCapacitySleepDuration time.Duration `mapstructure:"full_capacity_sleep_duration"`
	SendWorkerCount           int           `mapstructure:"send_worker_count"`
	MessageCost               int           `mapstructure:"message_cost"`
}

type SmsGateway struct {
	user usersrv.IUserService
	sms  smssrv.ISmsService
	cfg  Config
}

func NewSmsGateway(
	cfg Config,
	user usersrv.IUserService,
	sms smssrv.ISmsService,
) *SmsGateway {
	return &SmsGateway{
		cfg:  cfg,
		user: user,
		sms:  sms,
	}
}

func (s *SmsGateway) EnqueueWorker(ctx context.Context) error {
	enqueued, err := s.sms.EnqueueEarliest(ctx, s.cfg.EnqueueCount)
	if err != nil {
		logPkg.Error(err, "Failed to enqueue sms")

		if errors.Is(err, smsmodels.InvalidQueueError) {
			return err
		}

		if errors.Is(err, smsmodels.NoCapacityInQueueError) {
			return err
		}

		return nil
	}
	logPkg.Debug("Enqueued sms: %d", enqueued)

	return nil
}

func (s *SmsGateway) StartEnqueueWorker(ctx context.Context) {
	go func() {
		var stopErr error
		for {
			stopErr = ctx.Err()
			if stopErr != nil {
				logPkg.Error(stopErr, "Context error")
				break
			}

			stopErr = s.EnqueueWorker(ctx)
			if stopErr != nil {
				if errors.Is(stopErr, smsmodels.NoCapacityInQueueError) {
					time.Sleep(s.cfg.FullCapacitySleepDuration)
					continue
				}

				break
			}
		}

		logPkg.Error(stopErr, "Enqueue worker is shutting down")
	}()
}

func (s *SmsGateway) SendWorker(ctx context.Context) error {
	msg, err := s.sms.SendFromQueue(ctx)
	if err != nil {
		logPkg.Error(err, "Failed to enqueue sms")

		if errors.Is(err, smsmodels.InvalidQueueError) {
			return err
		}

		if errors.Is(err, smsmodels.MessageNotExistError) {
			return nil
		}

		return nil
	}

	if msg.Status == smsmodels.StatusFailed {
		if err := s.user.IncreaseUserBalance(ctx, msg.UserId, int64(msg.Cost)); err != nil {
			logPkg.Error(err, "Failed to increase user balance")
		}
	}

	return nil
}

func (s *SmsGateway) StartSendWorkers(ctx context.Context) {
	for range s.cfg.SendWorkerCount {
		go func() {
			var stopErr error
			for {
				stopErr = ctx.Err()
				if stopErr != nil {
					logPkg.Error(stopErr, "Context error")
					break
				}

				stopErr = s.SendWorker(ctx)
				if stopErr != nil {
					break
				}
			}

			logPkg.Error(stopErr, "Send worker is shutting down")
		}()
	}
}

func (s *SmsGateway) CreateUser(ctx context.Context, user usermodels.User) (usermodels.User, error) {
	res, err := s.user.CreateUser(ctx, user)
	if err != nil {
		logPkg.Error(err, "Failed to create user")
		return usermodels.User{}, err
	}

	return res, nil
}

func (s *SmsGateway) GetUser(ctx context.Context, userId string) (usermodels.User, error) {
	res, err := s.user.GetUser(ctx, userId)
	if err != nil {
		logPkg.Error(err, "Failed to get user")
		return usermodels.User{}, err
	}

	return res, nil
}

func (s *SmsGateway) GetUserMessages(ctx context.Context, userId string) ([]smsmodels.Sms, error) {
	userMsgs, err := s.sms.GetUserSms(ctx, userId)
	if err != nil {
		logPkg.Error(err, "Failed to get user messages")
		return nil, err
	}

	return userMsgs, nil
}

func (s *SmsGateway) SendSingleMessage(ctx context.Context, userId string, sms smsmodels.Sms) error {
	user, getUserErr := s.GetUser(ctx, userId)
	if getUserErr != nil {
		return getUserErr
	}

	totalCost := int64(s.cfg.MessageCost)

	if user.Balance < totalCost {
		return usermodels.InsufficientBalanceError
	}

	msg := smsmodels.Sms{
		Content:  sms.Content,
		Receiver: sms.Receiver,
		Cost:     s.cfg.MessageCost,
	}
	if decreaseErr := s.user.DecreaseUserBalance(ctx, userId, totalCost); decreaseErr != nil {
		logPkg.Error(decreaseErr, "Failed to decrease user balance")
		return decreaseErr
	}

	if scheduleErr := s.sms.ScheduleSms(ctx, userId, []smsmodels.Sms{msg}); scheduleErr != nil {
		logPkg.Error(scheduleErr, "Failed to schedule sms")

		if increaseErr := s.user.IncreaseUserBalance(ctx, userId, totalCost); increaseErr != nil {
			logPkg.Error(increaseErr, "Failed to increase user balance")
		}

		return scheduleErr
	}

	return nil
}

func (s *SmsGateway) SendBulkMessage(ctx context.Context, userId string, sms []smsmodels.Sms) error {
	user, getUserErr := s.GetUser(ctx, userId)
	if getUserErr != nil {
		return getUserErr
	}

	totalCost := int64(s.cfg.MessageCost * len(sms))

	if user.Balance < totalCost {
		return usermodels.InsufficientBalanceError
	}

	msgs := make([]smsmodels.Sms, len(sms))
	for i := range sms {
		msgs[i] = smsmodels.Sms{
			Content:  sms[i].Content,
			Receiver: sms[i].Receiver,
			Cost:     s.cfg.MessageCost,
		}
	}
	if decreaseErr := s.user.DecreaseUserBalance(ctx, userId, totalCost); decreaseErr != nil {
		logPkg.Error(decreaseErr, "Failed to decrease user balance")
		return decreaseErr
	}

	if scheduleErr := s.sms.ScheduleSms(ctx, userId, msgs); scheduleErr != nil {
		logPkg.Error(scheduleErr, "Failed to schedule sms")

		if increaseErr := s.user.IncreaseUserBalance(ctx, userId, totalCost); increaseErr != nil {
			logPkg.Error(increaseErr, "Failed to increase user balance")
		}

		return scheduleErr
	}

	return nil
}
