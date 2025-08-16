package services

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/repositories"
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
	for i := range msgs {
		msgs[i].UserId = userId
		msgs[i].Status = models.StatusScheduled
	}

	if err := s.smsRepo.CreateScheduleMessages(ctx, msgs); err != nil {
		return err
	}

	return nil
}

func (s *SmsService) GetUserSms(ctx context.Context, userId string, skip int, limit int, desc bool) ([]models.Sms, error) {
	msgs, err := s.smsRepo.GetMessagesByUserId(ctx, userId, skip, limit, desc)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *SmsService) EnqueueEarliest(ctx context.Context, count int) (int, error) {
	queueLen, err := s.smsQueue.GetLength(ctx)
	if err != nil {
		return 0, err
	}
	if queueLen+count > s.cfg.QueueCapacity {
		return 0, models.NoCapacityInQueueError
	}

	enqueuedMsgs, err := s.smsRepo.EnqueueMessages(ctx, count)
	if err != nil {
		return 0, err
	}

	if len(enqueuedMsgs) == 0 {
		return 0, nil
	}

	if err := s.smsQueue.Enqueue(ctx, enqueuedMsgs); err != nil {
		ids := make([]string, len(enqueuedMsgs))
		for i := range enqueuedMsgs {
			ids[i] = enqueuedMsgs[i].UserId
		}

		if scheduleErr := s.smsRepo.RescheduledMessages(ctx, ids); scheduleErr != nil {
			return 0, scheduleErr
		}
		return 0, err
	}

	return len(enqueuedMsgs), nil
}

func (s *SmsService) SetMessageAsFailed(ctx context.Context, id string) (models.Sms, error) {
	res, err := s.smsRepo.SetMessageAsFailed(ctx, id)
	if err != nil {
		return models.Sms{}, err
	}

	return res, nil
}

func (s *SmsService) SetMessageAsSent(ctx context.Context, id string) (models.Sms, error) {
	res, err := s.smsRepo.SetMessageAsSent(ctx, id)
	if err != nil {
		return models.Sms{}, err
	}

	return res, nil
}

func (s *SmsService) SendFromQueue(ctx context.Context) (models.Sms, error) {
	msg, err := s.smsQueue.Pop(ctx)
	if err != nil {
		return models.Sms{}, err
	}

	if err := s.smsSender.Send(ctx, msg); err != nil {
		return s.SetMessageAsFailed(ctx, msg.ID)
	}

	return s.SetMessageAsSent(ctx, msg.ID)
}
