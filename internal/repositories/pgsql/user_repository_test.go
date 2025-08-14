package pgsql_test

import (
	"context"
	"os"
	"testing"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/pgsql"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/mocks"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/repositories"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
	"github.com/stretchr/testify/assert"

	pkgPgSql "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/pgsql"
)

func initDB(smsQueue repositories.ISmsQueue) (*pkgPgSql.Connector, *pgsql.Repository, error) {
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

	repo, err := pgsql.NewRepository(conn, smsQueue)
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
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

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
		assert.NoError(t, err)
	})

	t.Run("should return EmptyNameError if user name is empty", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

		inputUser := models.User{
			Name:    "",
			Balance: 0,
		}

		actualUser, actualErr := repo.CreateUser(context.Background(), inputUser)

		assert.Error(t, actualErr)
		assert.Equal(t, models.EmptyNameError, actualErr)
		assert.Equal(t, "", actualUser.Name)
		assert.EqualValues(t, 0, actualUser.Balance)
		assert.Nil(t, actualUser.Entity)
		assert.Nil(t, actualUser.CreateDate)
		assert.Nil(t, actualUser.UpdateDate)

		err = cleanDB(conn)
		assert.NoError(t, err)
	})
}

func TestRepository_GetUser(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

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
		assert.NoError(t, err)
	})

	t.Run("should return UserNotExistError if user not exists", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

		actualUser, actualErr := repo.GetUser(context.Background(), "1")

		assert.Error(t, actualErr)
		assert.Equal(t, models.UserNotExistError, actualErr)
		assert.Equal(t, "", actualUser.Name)
		assert.EqualValues(t, 0, actualUser.Balance)
		assert.Nil(t, actualUser.Entity)
		assert.Nil(t, actualUser.CreateDate)
		assert.Nil(t, actualUser.UpdateDate)

		err = cleanDB(conn)
		assert.NoError(t, err)
	})
}

func TestRepository_UpdateUserBalance(t *testing.T) {
	t.Run("should update user balance", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

		inputUser := models.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		inputAmount := int64(100)

		ctx := context.Background()

		createdUser, err := repo.CreateUser(ctx, inputUser)
		assert.NoError(t, err)

		beforeUpdateUser, err := repo.GetUser(ctx, createdUser.ID)
		assert.NoError(t, err)

		actualErr := repo.UpdateUserBalance(ctx, createdUser.ID, inputAmount)
		assert.NoError(t, actualErr)

		afterUpdateUser, err := repo.GetUser(ctx, createdUser.ID)
		assert.NoError(t, err)

		assert.Equal(t, beforeUpdateUser.ID, afterUpdateUser.ID)
		assert.Equal(t, beforeUpdateUser.Name, afterUpdateUser.Name)
		assert.Equal(t, beforeUpdateUser.Balance+inputAmount, afterUpdateUser.Balance)
		assert.True(t, afterUpdateUser.CreatedAt.Equal(beforeUpdateUser.CreatedAt))
		assert.True(t, afterUpdateUser.UpdatedAt.After(beforeUpdateUser.UpdatedAt))

		err = cleanDB(conn)
		assert.NoError(t, err)
	})

	t.Run("should return UserNotExistError if user not exists", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

		ctx := context.Background()

		actualErr := repo.UpdateUserBalance(ctx, "1", 100)
		assert.Error(t, actualErr)
		assert.Equal(t, models.UserNotExistError, actualErr)

		err = cleanDB(conn)
		assert.NoError(t, err)
	})

	t.Run("should return InsufficientBalanceError if balance lt 0", func(t *testing.T) {
		conn, repo, err := initDB(mocks.NewMockISmsQueue(t))
		assert.NoError(t, err)

		inputUser := models.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		inputAmount := int64(-100)

		ctx := context.Background()

		createdUser, err := repo.CreateUser(ctx, inputUser)
		assert.NoError(t, err)

		actualErr := repo.UpdateUserBalance(ctx, createdUser.ID, inputAmount)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InsufficientBalanceError, actualErr)

		err = cleanDB(conn)
		assert.NoError(t, err)
	})
}
