package task

import (
	"context"
	"projects/LDmitryLD/task-service/task/internal/models"
	"projects/LDmitryLD/task-service/task/internal/modules/task/service"

	pbtask "github.com/LDmitryLD/protos2/gen/task"
	pbusers "github.com/LDmitryLD/protos2/gen/user"
)

type TaskServiceGRPC struct {
	taskService service.Tasker
	pbtask.UnimplementedTaskerServer
}

func NewTaskServiceGRPC(taskService service.Tasker) *TaskServiceGRPC {
	return &TaskServiceGRPC{taskService: taskService}
}

func (t *TaskServiceGRPC) Create(ctx context.Context, in *pbtask.CreateRequest) (*pbtask.CreateResponse, error) {
	id, err := t.taskService.Create(ctx, models.Task{TaskName: in.TaskName, Description: in.Description, UserId: int(in.UserId)})
	if err != nil {
		return nil, err
	}

	return &pbtask.CreateResponse{Id: uint32(id)}, nil
}

func (t *TaskServiceGRPC) Delete(ctx context.Context, in *pbtask.DeleteRequest) (*pbtask.DeleteResponse, error) {
	err := t.taskService.Delete(ctx, int(in.TaskId), int(in.UserId))
	if err != nil {
		return &pbtask.DeleteResponse{Success: false}, err
	}

	return &pbtask.DeleteResponse{Success: true}, nil
}

func (t *TaskServiceGRPC) List(ctx context.Context, in *pbtask.ListRequest) (*pbtask.ListResponse, error) {
	list, err := t.taskService.List(ctx, int(in.GetUserId()))
	if err != nil {
		return nil, err
	}

	tasks := make([]*pbusers.Task, len(list))
	for i, task := range list {
		tasks[i] = &pbusers.Task{
			TaskId:      uint32(task.TaskId),
			TaskName:    task.TaskName,
			Description: task.Description,
		}
	}

	return &pbtask.ListResponse{Tasks: tasks}, nil
}
