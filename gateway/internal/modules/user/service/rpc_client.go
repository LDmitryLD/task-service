package service

import (
	"context"
	"fmt"
	"projects/LDmitryLD/task-service/gateway/config"
	"projects/LDmitryLD/task-service/gateway/internal/models"

	pb "github.com/LDmitryLD/protos2/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Userer interface {
	Profile(ctx context.Context, id int) (models.User, error)
	Create(ctx context.Context, user models.User) (int, error)
}

type UserGRPClient struct {
	client pb.UsererClient
}

func NewUserGRPCClient(conf config.UserRPC) (Userer, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUsererClient(conn)
	return &UserGRPClient{client: client}, nil
}

func (u *UserGRPClient) Profile(ctx context.Context, id int) (models.User, error) {
	res, err := u.client.Profile(ctx, &pb.ProfileRequest{Id: uint32(id)})
	if err != nil {
		return models.User{}, err
	}

	tasks := make([]models.Task, len(res.Tasks))
	for i, task := range res.Tasks {
		tasks[i] = models.Task{
			TaskID:      int(task.TaskId),
			TaskName:    task.GetTaskName(),
			Description: task.GetDescription(),
		}
	}

	user := models.User{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Tasks:     tasks,
	}

	return user, nil
}

func (u *UserGRPClient) Create(ctx context.Context, user models.User) (int, error) {
	res, err := u.client.Create(ctx, &pb.CreateRequest{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email})
	if err != nil {
		return 0, err
	}

	return int(res.GetId()), nil
}
