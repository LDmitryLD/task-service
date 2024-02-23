package service

import (
	"context"
	"fmt"
	"projects/LDmitryLD/task-service/task/config"

	userpb "github.com/LDmitryLD/protos2/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Userer interface {
	Exists(ctx context.Context, userId int) bool
}

type UserGRPCClient struct {
	client userpb.UsererClient
}

func NewUserGRPCClient(conf config.UserRPC) (Userer, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := userpb.NewUsererClient(conn)
	return &UserGRPCClient{
		client: client,
	}, nil
}

func (t *UserGRPCClient) Exists(ctx context.Context, userId int) bool {
	req, err := t.client.Exists(ctx, &userpb.ExistsRequest{Id: uint32(userId)})
	if err != nil {
		return false
	}

	return req.Exists
}
