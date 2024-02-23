package service

import (
	"context"
	"fmt"
	"projects/LDmitryLD/task-service/user/config"

	taskpb "github.com/LDmitryLD/protos2/gen/task"
	userpb "github.com/LDmitryLD/protos2/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Tasker interface {
	List(ctx context.Context, userID int) ([]*userpb.Task, error)
}

type TaskGRPCClient struct {
	client taskpb.TaskerClient
}

func NewTaskGRPCClient(conf config.TaskRPC) (Tasker, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := taskpb.NewTaskerClient(conn)
	return &TaskGRPCClient{
		client: client,
	}, nil
}

func (t *TaskGRPCClient) List(ctx context.Context, userID int) ([]*userpb.Task, error) {
	res, err := t.client.List(ctx, &taskpb.ListRequest{UserId: uint32(userID)})
	if err != nil {
		return nil, err
	}

	return res.Tasks, nil
}
