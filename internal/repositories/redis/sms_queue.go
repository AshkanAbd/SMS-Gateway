package redis

import (
	"context"
	"errors"

	"github.com/AshkanAbd/arvancloud_sms_gateway/common"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	"github.com/redis/go-redis/v9"
)

const (
	wrongTypeError = "WRONGTYPE Operation against a key holding the wrong kind of value"
)

func (r *Repository) Enqueue(ctx context.Context, msg []models.Sms) error {
	strValues := make([]any, len(msg))
	for i := range msg {
		strValues[i] = common.ValueToJSON(msg[i])
	}

	_, err := r.queueClient.LPush(ctx, r.cfg.QueueName, strValues...).Result()
	if err != nil {
		if err.Error() == wrongTypeError {
			return models.InvalidQueueError
		}

		return err
	}

	return nil
}

func (r *Repository) GetLength(ctx context.Context) (int, error) {
	count, err := r.queueClient.LLen(ctx, r.cfg.QueueName).Result()
	if err != nil {
		if err.Error() == wrongTypeError {
			return 0, models.InvalidQueueError
		}

		return 0, err
	}

	return int(count), err
}

func (r *Repository) Pop(ctx context.Context) (models.Sms, error) {
	res, err := r.queueClient.BRPop(ctx, r.cfg.QueueTimeout, r.cfg.QueueName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.Sms{}, models.EmptyQueueError
		}
		if err.Error() == wrongTypeError {
			return models.Sms{}, models.InvalidQueueError
		}

		return models.Sms{}, err
	}

	s := models.Sms{}
	if err := common.JSONToValue[models.Sms](res[1], &s); err != nil {
		return models.Sms{}, err
	}

	return s, nil
}
