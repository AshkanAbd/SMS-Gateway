package config

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/services"
	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
)

type AppConfig struct {
	SmsGatewayConfig services.SmsGatewayConfig `mapstructure:"sms_gateway"`
	PgSQLConfig      pkgPgSql.Config           `mapstructure:"pgsql"`
	LogLevel         string                    `mapstructure:"log_level"`
}
