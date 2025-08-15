package pgsql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
)

type Repository struct {
	conn *gorm.DB
}

func NewRepository(pgsqlConn *pkgPgSql.Connector) (*Repository, error) {
	cfg := &gorm.Config{}

	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pgsqlConn.GetConnection(),
	}), cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		conn: conn,
	}, nil
}
