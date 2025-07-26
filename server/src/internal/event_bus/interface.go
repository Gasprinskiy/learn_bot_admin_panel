package event_bus

import (
	"context"
	"learn_bot_admin_panel/internal/entity/profile"
)

type UUIDChanel interface {
	Subscribe(ctx context.Context, UUID string) (profile.User, error)
	Publish(ctx context.Context, UUID string, user profile.User) error
}
