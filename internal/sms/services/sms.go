package services

import (
	"context"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/repositories"
)

type SmsServiceConfig struct {
	WorkerCount int `mapstructure:"worker_count"`
	PeekCount   int `mapstructure:"peek_count"`
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
