package services

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/repositories"
)

type SmsGatewayConfig struct {
	WorkerCount int `mapstructure:"worker_count"`
	PeekCount   int `mapstructure:"peek_count"`
}

type SmsGateway struct {
	smsRepo   repositories.ISmsRepository
	smsSender repositories.ISmsSender
	cfg       SmsGatewayConfig
}

func NewSmsGateway(
	cfg SmsGatewayConfig,
	smsRepo repositories.ISmsRepository,
	smsSender repositories.ISmsSender,
) *SmsGateway {
	return &SmsGateway{
		cfg:       cfg,
		smsRepo:   smsRepo,
		smsSender: smsSender,
	}
}

func (s *SmsGateway) StartWorkers() {

}

func (s *SmsGateway) Enqueue(msgs []models.Message) error {
	return nil
}
