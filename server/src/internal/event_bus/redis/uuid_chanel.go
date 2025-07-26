package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/event_bus"

	"github.com/redis/go-redis/v9"
)

type uuIDChanel struct {
	rdb *redis.Client
}

func NewUUIDChanel(rdb *redis.Client) event_bus.UUIDChanel {
	return &uuIDChanel{rdb}
}

func (r *uuIDChanel) getKey(UUID string) string {
	return fmt.Sprintf("uuid_chanel:%s", UUID)
}

func (r *uuIDChanel) Subscribe(ctx context.Context, UUID string) (profile.User, error) {
	var data profile.User

	sub := r.rdb.Subscribe(ctx, r.getKey(UUID))
	defer sub.Close()

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
		return data, err
	}

	return data, nil
}

func (r *uuIDChanel) Publish(ctx context.Context, UUID string, user profile.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.rdb.Publish(ctx, r.getKey(UUID), data).Err()
}
