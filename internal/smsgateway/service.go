package smsgateway

import (
	"context"
	"errors"
	"time"

	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	smssrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	usermodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
	usersrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/services"
	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

type Config struct {
	EnqueueCount              int           `mapstructure:"enqueue_count"`
	FullCapacitySleepDuration time.Duration `mapstructure:"full_capacity_sleep_duration"`
	EmptyEnqueueSleepDuration time.Duration `mapstructure:"empty_enqueue_sleep_duration"`
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

func (s *SmsGateway) EnqueueWorker(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "enqueue worker context canceled")
		return 0, err
	}

	newCtx := context.Background()
	enqueued, err := s.sms.EnqueueEarliest(newCtx, s.cfg.EnqueueCount)
	if err != nil {
		pkgLog.Error(err, "failed to enqueue sms")

		if errors.Is(err, smsmodels.InvalidQueueError) {
			return 0, err
		}

		if errors.Is(err, smsmodels.NoCapacityInQueueError) {
			return 0, err
		}

		return 0, nil
	}
	pkgLog.Debug("enqueued sms: %d", enqueued)

	return enqueued, nil
}

func (s *SmsGateway) StartEnqueueWorker(ctx context.Context) error {
	pkgLog.Debug("starting enqueue worker...")
	var stopErr error
	for {
		stopErr = ctx.Err()
		if stopErr != nil {
			pkgLog.Error(stopErr, "enqueue worker context canceled")
			break
		}

		enqueued, stopErr := s.EnqueueWorker(ctx)
		if stopErr != nil {
			if errors.Is(stopErr, smsmodels.NoCapacityInQueueError) {
				time.Sleep(s.cfg.FullCapacitySleepDuration)
				continue
			}

			break
		}
		if enqueued == 0 {
			time.Sleep(s.cfg.EmptyEnqueueSleepDuration)
		}
	}

	pkgLog.Error(stopErr, "enqueue worker shutdown successfully")
	return stopErr
}

func (s *SmsGateway) SendWorker(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "send worker context canceled")
		return err
	}

	newCtx := context.Background()
	msg, err := s.sms.SendFromQueue(newCtx)
	if err != nil {
		if errors.Is(err, smsmodels.MessageNotExistError) || errors.Is(err, smsmodels.EmptyQueueError) {
			return nil
		}
		pkgLog.Error(err, "failed to enqueue sms")

		if errors.Is(err, smsmodels.InvalidQueueError) {
			return err
		}

		return nil
	}

	if msg.Status == smsmodels.StatusFailed {
		if _, err := s.user.IncreaseUserBalance(newCtx, msg.UserId, int64(msg.Cost)); err != nil {
			pkgLog.Error(err, "failed to increase user balance")
		}
	}

	return nil
}

func (s *SmsGateway) StartSendWorkers(ctx context.Context) error {
	pkgLog.Debug("starting send worker...")
	var stopErr error
	for {
		stopErr = ctx.Err()
		if stopErr != nil {
			pkgLog.Error(stopErr, "send worker context canceled")
			break
		}

		stopErr = s.SendWorker(ctx)
		if stopErr != nil {
			break
		}
	}

	pkgLog.Error(stopErr, "send worker shutdown successfully")
	return stopErr
}

func (s *SmsGateway) CreateUser(ctx context.Context, user usermodels.User) (usermodels.User, error) {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "create user context canceled")
		return usermodels.User{}, err
	}

	newCtx := context.Background()
	res, err := s.user.CreateUser(newCtx, user)
	if err != nil {
		pkgLog.Error(err, "failed to create user")
		return usermodels.User{}, err
	}

	return res, nil
}

func (s *SmsGateway) GetUser(ctx context.Context, userId string) (usermodels.User, error) {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "get user context canceled")
		return usermodels.User{}, err
	}

	newCtx := context.Background()
	res, err := s.user.GetUser(newCtx, userId)
	if err != nil {
		pkgLog.Error(err, "failed to get user")
		return usermodels.User{}, err
	}

	return res, nil
}

func (s *SmsGateway) GetUserMessages(ctx context.Context, userId string, skip int, limit int, desc bool) ([]smsmodels.Sms, error) {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "get user context canceled")
		return nil, err
	}

	newCtx := context.Background()
	userMsgs, err := s.sms.GetUserSms(newCtx, userId, skip, limit, desc)
	if err != nil {
		pkgLog.Error(err, "failed to get user messages")
		return nil, err
	}

	return userMsgs, nil
}

func (s *SmsGateway) SendSingleMessage(ctx context.Context, userId string, sms smsmodels.Sms) error {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "send single message context canceled")
		return err
	}

	newCtx := context.Background()
	user, getUserErr := s.GetUser(newCtx, userId)
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
	if _, decreaseErr := s.user.DecreaseUserBalance(newCtx, userId, totalCost); decreaseErr != nil {
		pkgLog.Error(decreaseErr, "failed to decrease user balance")
		return decreaseErr
	}

	if scheduleErr := s.sms.ScheduleSms(newCtx, userId, []smsmodels.Sms{msg}); scheduleErr != nil {
		pkgLog.Error(scheduleErr, "failed to schedule sms")

		if _, increaseErr := s.user.IncreaseUserBalance(newCtx, userId, totalCost); increaseErr != nil {
			pkgLog.Error(increaseErr, "failed to increase user balance")
		}

		return scheduleErr
	}

	return nil
}

func (s *SmsGateway) SendBulkMessage(ctx context.Context, userId string, sms []smsmodels.Sms) error {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "send bulk message context canceled")
		return err
	}

	newCtx := context.Background()
	user, getUserErr := s.GetUser(newCtx, userId)
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
	if _, decreaseErr := s.user.DecreaseUserBalance(newCtx, userId, totalCost); decreaseErr != nil {
		pkgLog.Error(decreaseErr, "failed to decrease user balance")
		return decreaseErr
	}

	if scheduleErr := s.sms.ScheduleSms(newCtx, userId, msgs); scheduleErr != nil {
		pkgLog.Error(scheduleErr, "failed to schedule sms")

		if _, increaseErr := s.user.IncreaseUserBalance(newCtx, userId, totalCost); increaseErr != nil {
			pkgLog.Error(increaseErr, "failed to increase user balance")
		}

		return scheduleErr
	}

	return nil
}

func (s *SmsGateway) IncreaseUserBalance(ctx context.Context, userId string, amount int64) (int64, error) {
	if err := ctx.Err(); err != nil {
		pkgLog.Error(err, "increase user context canceled")
		return 0, err
	}

	newCtx := context.Background()
	newBalance, err := s.user.IncreaseUserBalance(newCtx, userId, amount)
	if err != nil {
		pkgLog.Error(err, "failed to increase user balance")
		return 0, err
	}

	return newBalance, nil
}
