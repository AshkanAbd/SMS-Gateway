package pgsql

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
)

type Repository struct {
	conn     *gorm.DB
	smsQueue repositories.ISmsQueue
}

func NewRepository(
	pgsqlConn *pkgPgSql.Connector,
	smsQueue repositories.ISmsQueue,
) (*Repository, error) {
	cfg := &gorm.Config{}

	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pgsqlConn.GetConnection(),
	}), cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		conn:     conn,
		smsQueue: smsQueue,
	}, nil
}
