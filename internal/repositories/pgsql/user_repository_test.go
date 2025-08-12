package pgsql_test

import (
	"context"
	"os"
	"testing"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/pgsql"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
	"github.com/stretchr/testify/assert"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
)

func initDB() (*pkgPgSql.Connector, *pgsql.Repository, error) {
	cfg := pkgPgSql.Config{
		DSN:            os.Getenv("PGSQL_DSN"),
		MigrationsPath: "file://../../../migrations/pgsql",
	}
	conn, err := pkgPgSql.NewConnector(cfg)

	if err != nil {
		return nil, nil, err
	}

	if err := conn.Migrate(); err != nil {
		return nil, nil, err
	}

	repo, err := pgsql.NewRepository(conn)
	if err != nil {
		return nil, nil, err
	}

	return conn, repo, nil
}

func cleanDB(conn *pkgPgSql.Connector) error {
	if err := conn.Clear(); err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func TestRepository_CreateUser(t *testing.T) {
	t.Run("should create user", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		inputUser := models.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}

		actualUser, actualErr := repo.CreateUser(context.Background(), inputUser)

		assert.NoError(t, actualErr)
		assert.Equal(t, inputUser.Name, actualUser.Name)
		assert.Equal(t, inputUser.Balance, actualUser.Balance)
		assert.NotNil(t, actualUser.Entity)
		assert.NotEqual(t, 0, actualUser.ID)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})

	t.Run("should return error if user name is empty", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		inputUser := models.User{
			Name:    "",
			Balance: 0,
		}

		actualUser, actualErr := repo.CreateUser(context.Background(), inputUser)

		assert.Error(t, actualErr)
		assert.Equal(t, "pq: new row for relation \"users\" violates check constraint \"users_name_check\"", actualErr.Error())
		assert.Equal(t, "", actualUser.Name)
		assert.EqualValues(t, 0, actualUser.Balance)
		assert.Nil(t, actualUser.Entity)
		assert.Nil(t, actualUser.CreateDate)
		assert.Nil(t, actualUser.UpdateDate)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})
}

func TestRepository_GetUser(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		inputUser := models.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}

		ctx := context.Background()

		expectedUser, err := repo.CreateUser(ctx, inputUser)
		assert.NoError(t, err)

		actualUser, actualErr := repo.GetUser(ctx, expectedUser.ID)

		assert.NoError(t, actualErr)
		assert.NotNil(t, actualUser.Entity)
		assert.NotNil(t, actualUser.CreateDate)
		assert.NotNil(t, actualUser.UpdateDate)
		assert.Equal(t, expectedUser.Name, actualUser.Name)
		assert.Equal(t, expectedUser.Balance, actualUser.Balance)
		assert.Equal(t, expectedUser.ID, actualUser.ID)
		assert.NotNil(t, actualUser.CreatedAt)
		assert.NotNil(t, actualUser.UpdatedAt)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})

	t.Run("should return error if user not exists", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		actualUser, actualErr := repo.GetUser(context.Background(), "1")

		assert.Error(t, actualErr)
		assert.Equal(t, models.UserNotExistError, actualErr)
		assert.Equal(t, "", actualUser.Name)
		assert.EqualValues(t, 0, actualUser.Balance)
		assert.Nil(t, actualUser.Entity)
		assert.Nil(t, actualUser.CreateDate)
		assert.Nil(t, actualUser.UpdateDate)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})
}
