package config

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/cmd/http/config"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/redis"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/smsgateway"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
	pkgRedis "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/redis"
)

type AppConfig struct {
	HttpConfig       config.HTTPConfig         `mapstructure:"http"`
	SmsServiceConfig services.SmsServiceConfig `mapstructure:"sms_service"`
	PgSQLConfig      pkgPgSql.Config           `mapstructure:"pgsql"`
	RedisConfig      pkgRedis.Config           `mapstructure:"redis"`
	RedisRepoConfig  redis.Config              `mapstructure:"redis_repo"`
	SmsGatewayConfig smsgateway.Config         `mapstructure:"sms_gateway"`
	LogLevel         string                    `mapstructure:"log_level"`
	SendWorkerCount  int                       `mapstructure:"send_worker_count"`
}
