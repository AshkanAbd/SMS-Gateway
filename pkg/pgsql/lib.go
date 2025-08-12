package pgsql

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	DSN            string `mapstructure:"dsn"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

type Connector struct {
	conn *sql.DB
	cfg  Config
}

func NewConnector(cfg Config) (*Connector, error) {
	connection, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &Connector{
		conn: connection,
		cfg:  cfg,
	}, nil
}

func (c *Connector) Migrate() error {
	driver, err := postgres.WithInstance(c.conn, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		c.cfg.MigrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (c *Connector) Clear() error {
	driver, err := postgres.WithInstance(c.conn, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		c.cfg.MigrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return err
	}

	err = migrator.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (c *Connector) Close() error {
	return c.conn.Close()
}

func (c *Connector) GetConnection() *sql.DB {
	return c.conn
}
