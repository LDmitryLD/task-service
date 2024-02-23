package storage

import (
	"context"
	"projects/LDmitryLD/task-service/task/internal/db/adapter"
	"projects/LDmitryLD/task-service/task/internal/models"
)

type TaskStorager interface {
	List(ctx context.Context, userId int) ([]models.Task, error)
	Create(ctx context.Context, task models.Task) (int, error)
	Delete(ctx context.Context, taskId int, userId int) error
}

type TaskStorage struct {
	adapter adapter.SQLAdapterer
}

func NewTaskStorage(sqlAdapter adapter.SQLAdapterer) TaskStorager {
	return &TaskStorage{
		adapter: sqlAdapter,
	}
}

func (t *TaskStorage) List(ctx context.Context, userId int) ([]models.Task, error) {
	return t.adapter.List(userId)
}

func (t *TaskStorage) Create(ctx context.Context, task models.Task) (int, error) {
	return t.adapter.Create(task)
}

func (t *TaskStorage) Delete(ctx context.Context, taskId int, userId int) error {
	return t.adapter.Delete(ctx, userId, taskId)
}
