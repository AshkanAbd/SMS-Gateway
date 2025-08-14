package services

import (
	"context"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/repositories"
)

type SmsServiceConfig struct {
	EnqueueInterval time.Duration `mapstructure:"enqueue_interval"`
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

func (s *SmsService) GetUserSms(ctx context.Context, userId string) ([]models.Sms, error) {
	msgs, err := s.smsRepo.GetMessagesByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *SmsService) EnqueueEarliest(ctx context.Context, count int) (int, error) {
	enqueuedCount := 0
	for range count {
		newEnqueued, err := s.smsRepo.EnqueueEarliestMessage(ctx)
		if err != nil {
			return 0, err
		}
		enqueuedCount += newEnqueued
	}

	return enqueuedCount, nil
}
