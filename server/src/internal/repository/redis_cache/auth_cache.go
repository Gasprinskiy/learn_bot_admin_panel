package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/tools/genredis"
	"strconv"

	"time"

	"github.com/redis/go-redis/v9"
)

type authCache struct {
	db  *redis.Client
	ttl time.Duration
}

func NewAuthCache(db *redis.Client, ttl time.Duration) repository.AuthCache {
	return &authCache{db, ttl}
}

func (r *authCache) SetTempUserData(ctx context.Context, tempKey string, user profile.User) error {
	byteData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.db.Set(ctx, tempKey, byteData, r.ttl).Err()
}

func (r *authCache) GetTempUserData(ctx context.Context, tempKey string) (profile.User, error) {
	return genredis.GetStruct[profile.User](ctx, r.db, tempKey)
}

func (r *authCache) SetDeletedUser(ctx context.Context, userID int) error {
	return r.db.Set(ctx, fmt.Sprintf("%d", userID), true, r.ttl*time.Hour).Err()
}

func (r *authCache) InUserDeletedCache(ctx context.Context, userID int) (bool, error) {
	result, err := r.db.Get(ctx, fmt.Sprintf("%d", userID)).Result()
	switch err {
	case nil:
		return strconv.ParseBool(result)

	case redis.Nil:
		return false, nil

	default:
		return false, err
	}
}
