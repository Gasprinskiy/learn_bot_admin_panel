package grpc_client

import (
	"context"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/repository/grpc_client/proto/notify_message"

	"google.golang.org/grpc"
)

type notifyMessageGRPCRepository struct {
	client notify_message.BotServiceClient
}

// NewAccessnodemonitGRPCRepository репозиторий accessnodemonit
func NewnotifyMessageGRPCRepository(conn *grpc.ClientConn) repository.NotifyMessage {
	return &notifyMessageGRPCRepository{
		notify_message.NewBotServiceClient(conn),
	}
}

func (r *notifyMessageGRPCRepository) SendInviteLink(ctx context.Context, TGID int64) (bool, error) {
	result, err := r.client.SendInviteLink(ctx, &notify_message.SendInviteLinkRequest{
		TgId: TGID,
	})
	if err != nil {
		return false, err
	}

	return result.Sent, nil
}
