package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/common"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/repositories/redis"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/stretchr/testify/assert"

	pkgRedis "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/redis"
)

const (
	queueName    = "test"
	queueDB      = 0
	queueTimeout = time.Second
)

func initRedis() (*pkgRedis.Connector, *redis.Repository, error) {
	connectorCfg := pkgRedis.Config{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	conn := pkgRedis.NewConnector(connectorCfg)

	repoCfg := redis.Config{
		QueueDB:      queueDB,
		QueueName:    queueName,
		QueueTimeout: queueTimeout,
	}

	repo := redis.NewRepository(repoCfg, conn)
	return conn, repo, nil
}

func cleanupRedis(conn *pkgRedis.Connector) error {
	if err := conn.GetClient(queueDB).FlushDB(context.Background()).Err(); err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func TestRepository_Enqueue(t *testing.T) {
	t.Run("should enqueue message", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()
		expectedMsg := models.Sms{
			UserId:   "1",
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}

		actualErr := repo.Enqueue(ctx, expectedMsg)
		assert.NoError(t, actualErr)

		res, err := conn.GetClient(queueDB).LRange(ctx, queueName, 0, -1).Result()
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		actualMsg := models.Sms{}
		err = common.JSONToValue(res[0], &actualMsg)
		assert.NoError(t, err)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})

	t.Run("should return InvalidQueueError when key is not list", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()
		expectedMsg := models.Sms{
			UserId:   "1",
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}

		err = conn.GetClient(queueDB).Set(ctx, queueName, 0, -1).Err()
		assert.NoError(t, err)

		actualErr := repo.Enqueue(ctx, expectedMsg)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidQueueError, actualErr)
	})
}

func TestRepository_GetLength(t *testing.T) {
	t.Run("should return queue length", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()
		msgStr := common.ValueToJSON(&models.Sms{
			UserId:   "1",
			Content:  "Test Content",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		})

		err = conn.GetClient(queueDB).LPush(ctx, queueName, msgStr).Err()
		assert.NoError(t, err)

		actualLen, actualErr := repo.GetLength(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, 1, actualLen)

		err = conn.GetClient(queueDB).LPush(ctx, queueName, msgStr).Err()
		assert.NoError(t, err)

		actualLen, actualErr = repo.GetLength(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, 2, actualLen)

		err = conn.GetClient(queueDB).RPop(ctx, queueName).Err()
		assert.NoError(t, err)

		actualLen, actualErr = repo.GetLength(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, 1, actualLen)
	})

	t.Run("should return InvalidQueueError when key is not list", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()

		err = conn.GetClient(queueDB).Set(ctx, queueName, 0, -1).Err()
		assert.NoError(t, err)

		actualLen, actualErr := repo.GetLength(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidQueueError, actualErr)
		assert.Equal(t, 0, actualLen)
	})
}

func TestRepository_Pop(t *testing.T) {
	t.Run("should pop and return oldest message", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()

		expectedMsg := models.Sms{
			UserId:   "1",
			Content:  "Test Content 1",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}
		err = conn.GetClient(queueDB).LPush(ctx, queueName, common.ValueToJSON(&expectedMsg)).Err()
		assert.NoError(t, err)

		err = conn.GetClient(queueDB).LPush(ctx, queueName, common.ValueToJSON(&models.Sms{
			UserId:   "2",
			Content:  "Test Content 2",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		})).Err()
		assert.NoError(t, err)

		err = conn.GetClient(queueDB).LPush(ctx, queueName, common.ValueToJSON(&models.Sms{
			UserId:   "3",
			Content:  "Test Content 3",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		})).Err()
		assert.NoError(t, err)

		actualMsg, actualErr := repo.Pop(ctx)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})

	t.Run("should return InvalidQueueError when key is not list", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()

		err = conn.GetClient(queueDB).Set(ctx, queueName, 0, -1).Err()
		assert.NoError(t, err)

		_, actualErr := repo.Pop(ctx)
		assert.Error(t, actualErr)
		assert.Equal(t, models.InvalidQueueError, actualErr)
	})

	t.Run("should block until message is available when list is empty", func(t *testing.T) {
		conn, repo, err := initRedis()
		assert.NoError(t, err)

		defer func() {
			err = cleanupRedis(conn)
			assert.NoError(t, err)
		}()

		ctx := context.Background()

		expectedMsg := models.Sms{
			UserId:   "1",
			Content:  "Test Content 1",
			Receiver: "09123456789",
			Cost:     100,
			Status:   models.StatusEnqueued,
		}

		before := time.Now()
		go func() {
			time.Sleep(1 * time.Second)
			err = conn.GetClient(queueDB).LPush(ctx, queueName, common.ValueToJSON(&expectedMsg)).Err()
			assert.NoError(t, err)
		}()
		actualMsg, actualErr := repo.Pop(ctx)
		assert.GreaterOrEqual(t, time.Now().Sub(before), 500*time.Millisecond)
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedMsg.UserId, actualMsg.UserId)
		assert.Equal(t, expectedMsg.Content, actualMsg.Content)
		assert.Equal(t, expectedMsg.Receiver, actualMsg.Receiver)
		assert.Equal(t, expectedMsg.Cost, actualMsg.Cost)
		assert.Equal(t, expectedMsg.Status, actualMsg.Status)
	})
}
