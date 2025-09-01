package grpc_client

import (
	"context"
	"learn_bot_admin_panel/internal/entity/chanel_kicker"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/repository/grpc_client/proto/kicker"

	"google.golang.org/grpc"
)

type kickerRepo struct {
	client kicker.KickerServiceClient
}

func NewKicker(conn *grpc.ClientConn) repository.Kicker {
	return &kickerRepo{
		kicker.NewKickerServiceClient(conn),
	}
}

func (r *kickerRepo) KickUser(param chanel_kicker.KickUserParam) error {
	_, err := r.client.KickUsers(context.Background(), &kicker.KickUsersRequest{
		Params: []*kicker.KickUserParam{{
			TgId:     param.TgID,
			ReasonId: int64(param.Reason),
		}},
	})

	return err
}
