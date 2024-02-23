package service

import (
	"context"
	"projects/LDmitryLD/task-service/task/internal/models"
	"projects/LDmitryLD/task-service/task/internal/modules/task/storage"
	"projects/LDmitryLD/task-service/task/internal/modules/user/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Tasker interface {
	List(ctx context.Context, userId int) ([]models.Task, error)
	Create(ctx context.Context, task models.Task) (int, error)
	Delete(ctx context.Context, taskId int, userId int) error
}

type Task struct {
	user    service.Userer
	storage storage.TaskStorager
	logger  *logrus.Logger
}

func NewTasker(storage storage.TaskStorager, user service.Userer, logger *logrus.Logger) Tasker {
	return &Task{
		user:    user,
		storage: storage,
		logger:  logger,
	}
}
func (t *Task) List(ctx context.Context, userId int) ([]models.Task, error) {
	if !t.user.Exists(ctx, userId) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return t.storage.List(ctx, userId)
}

func (t *Task) Create(ctx context.Context, task models.Task) (int, error) {
	if !t.user.Exists(ctx, task.UserId) {
		return 0, status.Error(codes.NotFound, "user not found")
	}
	return t.storage.Create(ctx, task)
}

func (t *Task) Delete(ctx context.Context, taskId int, userId int) error {
	if !t.user.Exists(ctx, userId) {
		return status.Error(codes.NotFound, "user not found")
	}

	return t.storage.Delete(ctx, taskId, userId)
}
