package service

import (
	"context"
	"fmt"
	"projects/LDmitryLD/task-service/gateway/config"
	"projects/LDmitryLD/task-service/gateway/internal/models"

	pb "github.com/LDmitryLD/protos2/gen/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Tasker interface {
	Create(ctx context.Context, task models.TaskDTO) (int, error)
	List(ctx context.Context, userID int) ([]models.Task, error)
	Delete(ctx context.Context, userID, taskID int) (bool, error)
}

type TaskGRPCClient struct {
	client pb.TaskerClient
}

func NewTaskGRPCClient(conf config.TaskRPC) (Tasker, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTaskerClient(conn)
	return &TaskGRPCClient{client: client}, nil
}

func (t *TaskGRPCClient) Create(ctx context.Context, task models.TaskDTO) (int, error) {
	res, err := t.client.Create(ctx, &pb.CreateRequest{TaskName: task.TaskName, Description: task.Description, UserId: uint32(task.UserID)})
	if err != nil {
		return 0, err
	}

	return int(res.GetId()), nil
}

func (t *TaskGRPCClient) List(ctx context.Context, userID int) ([]models.Task, error) {
	res, err := t.client.List(ctx, &pb.ListRequest{UserId: uint32(userID)})
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, len(res.Tasks))
	for i, task := range res.Tasks {
		tasks[i] = models.Task{
			TaskID:      int(task.GetTaskId()),
			TaskName:    task.GetTaskName(),
			Description: task.GetDescription(),
		}
	}

	return tasks, nil
}

func (t *TaskGRPCClient) Delete(ctx context.Context, userID, taskID int) (bool, error) {
	res, err := t.client.Delete(ctx, &pb.DeleteRequest{UserId: uint32(userID), TaskId: uint32(taskID)})
	if err != nil {
		return false, err
	}

	return res.Success, nil
}
