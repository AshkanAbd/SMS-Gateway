package config

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/services"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/redis"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
	pkgRedis "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/redis"
)

type AppConfig struct {
	SmsServiceConfig services.SmsServiceConfig `mapstructure:"sms_service"`
	PgSQLConfig      pkgPgSql.Config           `mapstructure:"pgsql"`
	RedisConfig      pkgRedis.Config           `mapstructure:"redis"`
	RedisRepoConfig  redis.Config              `mapstructure:"redis_repo"`
	LogLevel         string                    `mapstructure:"log_level"`
}
