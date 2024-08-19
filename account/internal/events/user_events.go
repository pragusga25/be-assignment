package events

import (
	"context"
	"encoding/json"

	"pragusga/internal/domain"

	"github.com/go-redis/redis/v8"
)

const (
	UserCreatedChannel = "user:created"
)

type UserEventPublisher struct {
	redisClient *redis.Client
}

func NewUserEventPublisher(redisClient *redis.Client) *UserEventPublisher {
	return &UserEventPublisher{redisClient: redisClient}
}

func (p *UserEventPublisher) PublishUserCreated(ctx context.Context, user *domain.User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return p.redisClient.Publish(ctx, UserCreatedChannel, userData).Err()
}
