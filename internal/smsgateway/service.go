package smsgateway

import (
	"context"
	"errors"

	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	smssrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	usersrv "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/services"
	logPkg "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

type Config struct {
	EnqueueCount    int `mapstructure:"enqueue_count"`
	SendWorkerCount int `mapstructure:"send_worker_count"`
}

type SmsGateway struct {
	user usersrv.UserService
	sms  smssrv.SmsService
	cfg  Config
}

func NewSmsGateway(
	cfg Config,
	user usersrv.UserService,
	sms smssrv.SmsService,
) *SmsGateway {
	return &SmsGateway{
		cfg:  cfg,
		user: user,
		sms:  sms,
	}
}

func (s *SmsGateway) enqueueWorker(ctx context.Context) error {
	enqueued, err := s.sms.EnqueueEarliest(ctx, s.cfg.EnqueueCount)
	if err != nil {
		logPkg.Error(err, "Failed to enqueue sms")

		if errors.Is(err, smsmodels.InvalidQueueError) {
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

			stopErr = s.enqueueWorker(ctx)
			if stopErr != nil {
				break
			}
		}

		logPkg.Error(stopErr, "Enqueue worker is shutting down")
	}()
}

func (s *SmsGateway) sendWorker(ctx context.Context) error {
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

				stopErr = s.sendWorker(ctx)
				if stopErr != nil {
					break
				}
			}

			logPkg.Error(stopErr, "Send worker is shutting down")
		}()
	}
}
