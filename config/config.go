package config

import "github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/services"

type AppConfig struct {
	SmsGatewayConfig services.SmsGatewayConfig `mapstructure:"sms_gateway"`
	LogLevel         string                    `mapstructure:"log_level"`
}
