package redis

import (
	"context"
	"errors"
	"fmt"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"strconv"

	"github.com/DedAzaMarks/ABS/internal/domain"

	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type Cache struct {
	redis *redis.Client
}

func New(_ context.Context, redis *redis.Client) *Cache {
	return &Cache{redis: redis}
}

func (c *Cache) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	var pbUser pb.TgUser
	buf, err := c.redis.Get(ctx, strconv.Itoa(int(userID))).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, myerrors.ErrorUserNotFound
		}
		return nil, fmt.Errorf("error on getting user%w", err)
	}
	if err := proto.Unmarshal([]byte(buf), &pbUser); err != nil {
		return nil, fmt.Errorf("error on unmarshalling user: %w", err)
	}
	user := domain.PB2TGUser(&pbUser)
	return user, nil
}

func (c *Cache) SetUser(ctx context.Context, userID int64, user *domain.User) error {
	buf, err := proto.Marshal(domain.TGUser2PB(user))
	if err != nil {
		return fmt.Errorf("error on marshalling user: %w", err)
	}
	if err := c.redis.Set(ctx, strconv.Itoa(int(userID)), string(buf), 0).Err(); err != nil {
		return fmt.Errorf("error on setting user: %w", err)
	}
	return nil
}
