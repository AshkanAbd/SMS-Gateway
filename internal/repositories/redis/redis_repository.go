package redis

import (
	"time"

	"github.com/redis/go-redis/v9"

	pkgRedis "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/redis"
)

type Config struct {
	QueueDB      int           `mapstructure:"queue_db"`
	QueueName    string        `mapstructure:"queue_name"`
	QueueTimeout time.Duration `mapstructure:"queue_timeout"`
}

type Repository struct {
	queueClient *redis.Client
	cfg         Config
}

func NewRepository(cfg Config, conn *pkgRedis.Connector) *Repository {
	return &Repository{
		cfg:         cfg,
		queueClient: conn.GetClient(cfg.QueueDB),
	}
}
