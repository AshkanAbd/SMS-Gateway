package services

import (
	"context"
	"errors"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/repositories"

	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
	pkgMetrics "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/metrics"
)

type SmsServiceConfig struct {
	QueueCapacity int `mapstructure:"queue_capacity"`
}

type ISmsService interface {
	ScheduleSms(ctx context.Context, userId string, msgs []models.Sms) error
	GetUserSms(ctx context.Context, userId string, skip int, limit int, desc bool) ([]models.Sms, error)
	EnqueueEarliest(ctx context.Context, count int) (int, error)
	SetMessageAsFailed(ctx context.Context, id string) (models.Sms, error)
	SetMessageAsSent(ctx context.Context, id string) (models.Sms, error)
	SendFromQueue(ctx context.Context) (models.Sms, error)
}

type SmsService struct {
	smsRepo   repositories.ISmsRepository
	smsSender repositories.ISmsSender
	smsQueue  repositories.ISmsQueue
	cfg       SmsServiceConfig
}

func NewSmsService(
	cfg SmsServiceConfig,
	smsRepo repositories.ISmsRepository,
	smsSender repositories.ISmsSender,
	smsQueue repositories.ISmsQueue,
) *SmsService {
	return &SmsService{
		cfg:       cfg,
		smsRepo:   smsRepo,
		smsSender: smsSender,
		smsQueue:  smsQueue,
	}
}

func (s *SmsService) ScheduleSms(ctx context.Context, userId string, msgs []models.Sms) error {
	pkgLog.Debug("scheduling %d sms for user %s", len(msgs), userId)
	for i := range msgs {
		msgs[i].UserId = userId
		msgs[i].Status = models.StatusScheduled
	}

	if err := s.smsRepo.CreateScheduleMessages(ctx, msgs); err != nil {
		pkgLog.Error(err, "error creating scheduled sms for user %s", userId)
		return err
	}
	pkgMetrics.SmsStatusMetric.WithLabelValues("scheduled").Add(float64(len(msgs)))

	pkgLog.Debug("%d sms scheduled for user %s", len(msgs), userId)
	return nil
}

func (s *SmsService) GetUserSms(ctx context.Context, userId string, skip int, limit int, desc bool) ([]models.Sms, error) {
	pkgLog.Debug("getting sms for user %s", userId)
	msgs, err := s.smsRepo.GetMessagesByUserId(ctx, userId, skip, limit, desc)
	if err != nil {
		pkgLog.Error(err, "error getting sms for user %s", userId)
		return nil, err
	}

	pkgLog.Debug("%d sms retrieved for user %s", len(msgs), userId)
	return msgs, nil
}

func (s *SmsService) EnqueueEarliest(ctx context.Context, count int) (int, error) {
	pkgLog.Debug("enqueuing %d sms...", count)
	pkgLog.Debug("getting sms queue length")
	queueLen, err := s.smsQueue.GetLength(ctx)
	if err != nil {
		pkgLog.Error(err, "failed to retrieve sms queue length")
		return 0, err
	}
	pkgLog.Debug("sms queue len is %d", queueLen)
	if queueLen+count > s.cfg.QueueCapacity {
		pkgLog.Error(models.NoCapacityInQueueError, "queue length exceeds limit %d", s.cfg.QueueCapacity)
		return 0, models.NoCapacityInQueueError
	}

	pkgLog.Debug("enqueuing %d sms", count)
	enqueuedMsgs, err := s.smsRepo.EnqueueMessages(ctx, count)
	if err != nil {
		pkgLog.Error(err, "failed to enqueue sms")
		return 0, err
	}

	if len(enqueuedMsgs) == 0 {
		pkgLog.Debug("no enqueued sms")
		return 0, nil
	}
	pkgMetrics.SmsStatusMetric.WithLabelValues("enqueued").Add(float64(len(enqueuedMsgs)))

	pkgLog.Debug("adding %d sms to queue", len(enqueuedMsgs))
	if err := s.smsQueue.Enqueue(ctx, enqueuedMsgs); err != nil {
		pkgLog.Error(err, "failed to add sms to queue")
		ids := make([]string, len(enqueuedMsgs))
		for i := range enqueuedMsgs {
			ids[i] = enqueuedMsgs[i].UserId
		}

		pkgLog.Debug("rescheduling %d sms...", len(enqueuedMsgs))
		if scheduleErr := s.smsRepo.RescheduledMessages(ctx, ids); scheduleErr != nil {
			pkgLog.Error(scheduleErr, "failed to reschedule sms")
			return 0, scheduleErr
		}
		return 0, err
	}

	pkgLog.Debug("%d sms enqueued successfully", len(enqueuedMsgs))
	return len(enqueuedMsgs), nil
}

func (s *SmsService) SetMessageAsFailed(ctx context.Context, id string) (models.Sms, error) {
	pkgLog.Debug("setting message %s as failed", id)
	res, err := s.smsRepo.SetMessageAsFailed(ctx, id)
	if err != nil {
		pkgLog.Error(err, "failed to set message as failed")
		return models.Sms{}, err
	}

	pkgMetrics.SmsStatusMetric.WithLabelValues("failed").Inc()

	pkgLog.Debug("message %s set as failed", id)
	return res, nil
}

func (s *SmsService) SetMessageAsSent(ctx context.Context, id string) (models.Sms, error) {
	pkgLog.Debug("setting message %s as sent", id)
	res, err := s.smsRepo.SetMessageAsSent(ctx, id)
	if err != nil {
		pkgLog.Error(err, "failed to set message as sent")
		return models.Sms{}, err
	}

	pkgMetrics.SmsStatusMetric.WithLabelValues("sent").Inc()

	pkgLog.Debug("message %s set as sent", id)
	return res, nil
}

func (s *SmsService) SendFromQueue(ctx context.Context) (models.Sms, error) {
	pkgLog.Debug("poping message from queue")
	msg, err := s.smsQueue.Pop(ctx)
	if err != nil {
		if errors.Is(err, models.EmptyQueueError) {
			return models.Sms{}, err
		}
		pkgLog.Error(err, "failed to pop message from queue")
		return models.Sms{}, err
	}

	pkgLog.Debug("trying to send message %s to sms provider", msg.ID)
	if err := s.smsSender.Send(ctx, msg); err != nil {
		pkgLog.Error(err, "failed to send message %s to sms provider", msg.ID)
		return s.SetMessageAsFailed(ctx, msg.ID)
	}

	return s.SetMessageAsSent(ctx, msg.ID)
}
