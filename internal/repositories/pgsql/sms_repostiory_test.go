package pgsql_test

import (
	"context"
	"testing"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/sms/models"
	"github.com/stretchr/testify/assert"

	umodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
)

func TestRepository_CreateScheduleMessages(t *testing.T) {
	t.Run("should create messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		ctx := context.Background()

		tmpUser := umodels.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		createdUser, err := repo.CreateUser(ctx, tmpUser)
		assert.NoError(t, err)

		inputMsgs := []models.Sms{
			{
				UserId:   createdUser.ID,
				Content:  "Test Content",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		row := conn.GetConnection().QueryRow("select count(1) from messages")
		count := 0
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), count)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})

	t.Run("should return EmptyContentError when content is empty and not create any messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		ctx := context.Background()

		tmpUser := umodels.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		createdUser, err := repo.CreateUser(ctx, tmpUser)
		assert.NoError(t, err)

		inputMsgs := []models.Sms{
			{
				UserId:   createdUser.ID,
				Content:  "Test Content",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
			{
				UserId:   createdUser.ID,
				Content:  "",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.Error(t, err)
		assert.Equal(t, models.EmptyContentError, err)

		row := conn.GetConnection().QueryRow("select count(1) from messages")
		count := 0
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})

	t.Run("should return EmptyReceiverError when receiver is empty and not create any messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		ctx := context.Background()

		tmpUser := umodels.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		createdUser, err := repo.CreateUser(ctx, tmpUser)
		assert.NoError(t, err)

		inputMsgs := []models.Sms{
			{
				UserId:   createdUser.ID,
				Content:  "Test Content",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content",
				Receiver: "",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.Error(t, err)
		assert.Equal(t, models.EmptyReceiverError, err)

		row := conn.GetConnection().QueryRow("select count(1) from messages")
		count := 0
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		err = cleanDB(conn)
		assert.Nil(t, err)
	})
}

func TestRepository_GetMessagesByUserId(t *testing.T) {
	t.Run("should return user messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.Nil(t, err)

		ctx := context.Background()

		tmpUser := umodels.User{
			Name:    "AshkanAbd",
			Balance: 0,
		}
		createdUser, err := repo.CreateUser(ctx, tmpUser)
		assert.NoError(t, err)

		inputMsgs := []models.Sms{
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 1",
				Receiver: "09123456788",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 2",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		actualMsgs, actualErr := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, actualErr)
		assert.Equal(t, len(inputMsgs), len(actualMsgs))
		for i := 0; i < len(inputMsgs); i++ {
			assert.NotNil(t, actualMsgs[i].Entity)
			assert.NotNil(t, actualMsgs[i].CreateDate)
			assert.NotNil(t, actualMsgs[i].UpdateDate)
			assert.NotNil(t, actualMsgs[i].ID)
			assert.NotNil(t, actualMsgs[i].CreatedAt)
			assert.NotNil(t, actualMsgs[i].UpdatedAt)

			assert.Equal(t, inputMsgs[i].UserId, actualMsgs[i].UserId)
			assert.Equal(t, inputMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, inputMsgs[i].Receiver, actualMsgs[i].Receiver)
			assert.Equal(t, inputMsgs[i].Cost, actualMsgs[i].Cost)
			assert.Equal(t, inputMsgs[i].Status, actualMsgs[i].Status)
		}

		err = cleanDB(conn)
		assert.Nil(t, err)
	})
}
