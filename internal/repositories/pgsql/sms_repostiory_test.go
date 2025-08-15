package pgsql_test

import (
	"context"
	"slices"
	"testing"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/stretchr/testify/assert"

	umodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
)

func TestRepository_CreateScheduleMessages(t *testing.T) {
	t.Run("should create messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.NoError(t, err)

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
		assert.NoError(t, err)
	})

	t.Run("should return EmptyContentError when content is empty and not create any messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.NoError(t, err)

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
		assert.NoError(t, err)
	})

	t.Run("should return EmptyReceiverError when receiver is empty and not create any messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.NoError(t, err)

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
		assert.NoError(t, err)
	})
}

func TestRepository_GetMessagesByUserId(t *testing.T) {
	t.Run("should return user messages", func(t *testing.T) {
		conn, repo, err := initDB()
		assert.NoError(t, err)

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
		slices.Reverse(inputMsgs)
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
		assert.NoError(t, err)
	})
}

func TestRepository_SetMessageAsFailed(t *testing.T) {
	t.Run("should set message as failed", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
				Status:   models.StatusEnqueued,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		actualMsg, actualErr := repo.SetMessageAsFailed(ctx, userMsgs[0].ID)
		assert.NoError(t, actualErr)
		assert.Equal(t, userMsgs[0].ID, actualMsg.ID)
		assert.Equal(t, userMsgs[0].CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, userMsgs[0].UserId, actualMsg.UserId)
		assert.Equal(t, userMsgs[0].Content, actualMsg.Content)
		assert.Equal(t, userMsgs[0].Receiver, actualMsg.Receiver)
		assert.Equal(t, userMsgs[0].Cost, actualMsg.Cost)
		assert.Equal(t, models.StatusFailed, actualMsg.Status)
		assert.True(t, userMsgs[0].UpdatedAt.Before(actualMsg.UpdatedAt))
	})

	t.Run("should return MessageNotExistError when message status is not enqueue", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		_, actualErr := repo.SetMessageAsFailed(ctx, userMsgs[0].ID)
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
	})

	t.Run("should return MessageNotExistError when message not exists", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

		_, actualErr := repo.SetMessageAsFailed(ctx, "1")
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
	})
}

func TestRepository_SetMessageAsSent(t *testing.T) {
	t.Run("should set message as sent", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
				Status:   models.StatusEnqueued,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		actualMsg, actualErr := repo.SetMessageAsSent(ctx, userMsgs[0].ID)
		assert.NoError(t, actualErr)
		assert.Equal(t, userMsgs[0].ID, actualMsg.ID)
		assert.Equal(t, userMsgs[0].CreatedAt, actualMsg.CreatedAt)
		assert.Equal(t, userMsgs[0].UserId, actualMsg.UserId)
		assert.Equal(t, userMsgs[0].Content, actualMsg.Content)
		assert.Equal(t, userMsgs[0].Receiver, actualMsg.Receiver)
		assert.Equal(t, userMsgs[0].Cost, actualMsg.Cost)
		assert.Equal(t, models.StatusSent, actualMsg.Status)
		assert.True(t, userMsgs[0].UpdatedAt.Before(actualMsg.UpdatedAt))
	})

	t.Run("should return MessageNotExistError when message status is not enqueue", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		_, actualErr := repo.SetMessageAsSent(ctx, userMsgs[0].ID)
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
	})

	t.Run("should return MessageNotExistError when message not exists", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

		_, actualErr := repo.SetMessageAsSent(ctx, "1")
		assert.Error(t, actualErr)
		assert.Equal(t, models.MessageNotExistError, actualErr)
	})
}

func TestRepository_EnqueueMessages(t *testing.T) {
	t.Run("should enqueue scheduled message with given count", func(t *testing.T) {
		ctx := context.Background()
		countInput := 2

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
				Receiver: "09123456789",
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
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 3",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusScheduled,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		actualMsgs, actualErr := repo.EnqueueMessages(ctx, countInput)
		assert.NoError(t, actualErr)
		assert.Equal(t, countInput, len(actualMsgs))
		slices.Reverse(userMsgs)
		for i := 0; i < countInput; i++ {
			assert.Equal(t, userMsgs[i].ID, actualMsgs[i].ID)
			assert.Equal(t, userMsgs[i].CreatedAt, actualMsgs[i].CreatedAt)
			assert.Equal(t, userMsgs[i].UserId, actualMsgs[i].UserId)
			assert.Equal(t, userMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, userMsgs[i].Receiver, actualMsgs[i].Receiver)
			assert.Equal(t, userMsgs[i].Cost, actualMsgs[i].Cost)
			assert.True(t, userMsgs[i].UpdatedAt.Before(actualMsgs[i].UpdatedAt))
			assert.Equal(t, models.StatusEnqueued, actualMsgs[i].Status)
		}
	})

	t.Run("should not enqueue if no message was in schedule", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusEnqueued,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 2",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusSent,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 3",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusFailed,
			},
		}
		countInput := len(inputMsgs)

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		actualMsgs, actualErr := repo.EnqueueMessages(ctx, countInput)
		assert.NoError(t, actualErr)
		assert.Equal(t, 0, len(actualMsgs))
	})
}

func TestRepository_RescheduledMessages(t *testing.T) {
	t.Run("should reschedule messages with given id", func(t *testing.T) {
		ctx := context.Background()

		conn, repo, err := initDB()
		assert.NoError(t, err)

		defer func() {
			err = cleanDB(conn)
			assert.NoError(t, err)
		}()

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
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusSent,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 2",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusEnqueued,
			},
			{
				UserId:   createdUser.ID,
				Content:  "Test Content 3",
				Receiver: "09123456789",
				Cost:     100,
				Status:   models.StatusEnqueued,
			},
		}

		err = repo.CreateScheduleMessages(ctx, inputMsgs)
		assert.NoError(t, err)

		userMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(inputMsgs), len(userMsgs))

		actualErr := repo.RescheduledMessages(ctx, []string{
			userMsgs[0].ID,
			userMsgs[2].ID,
		})
		assert.NoError(t, actualErr)

		actualMsgs, err := repo.GetMessagesByUserId(ctx, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, len(actualMsgs), len(userMsgs))
		for i := 0; i < len(userMsgs); i++ {
			assert.Equal(t, userMsgs[i].ID, actualMsgs[i].ID)
			assert.Equal(t, userMsgs[i].CreatedAt, actualMsgs[i].CreatedAt)
			assert.Equal(t, userMsgs[i].UserId, actualMsgs[i].UserId)
			assert.Equal(t, userMsgs[i].Content, actualMsgs[i].Content)
			assert.Equal(t, userMsgs[i].Receiver, actualMsgs[i].Receiver)
			assert.Equal(t, userMsgs[i].Cost, actualMsgs[i].Cost)
			if i == 0 || i == 2 {
				assert.True(t, userMsgs[i].UpdatedAt.Before(actualMsgs[i].UpdatedAt))
				assert.Equal(t, models.StatusScheduled, actualMsgs[i].Status)
			} else {
				assert.Equal(t, userMsgs[i].UpdatedAt, actualMsgs[i].UpdatedAt)
				assert.Equal(t, userMsgs[i].Status, actualMsgs[i].Status)
			}
		}
	})
}
