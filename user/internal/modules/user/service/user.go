package service

import (
	"context"
	"projects/LDmitryLD/task-service/user/internal/models"
	"projects/LDmitryLD/task-service/user/internal/modules/task/service"
	"projects/LDmitryLD/task-service/user/internal/modules/user/storage"

	"go.uber.org/zap"
)

type Userer interface {
	Profile(ctx context.Context, id int) (models.User, error)
	Create(user models.User) (int, error)
	Exists(ctx context.Context, id int) error
}

type User struct {
	task    service.Tasker
	storage storage.UserStorager
	logger  *zap.Logger
}

func NewUser(storage storage.UserStorager, task service.Tasker, logger *zap.Logger) Userer {
	return &User{
		task:    task,
		storage: storage,
		logger:  logger,
	}
}

func (u *User) Profile(ctx context.Context, id int) (models.User, error) {
	user, err := u.storage.GetByID(id)
	if err != nil {
		u.logger.Error("user: User.Profile GetByID err", zap.Error(err))
		return models.User{}, err
	}

	tasks, err := u.task.List(ctx, id)
	if err != nil {
		u.logger.Error("user: User.Profile List err", zap.Error(err))
		return models.User{}, err
	}
	user.Tasks = tasks

	return user, nil
}

func (u *User) Create(user models.User) (int, error) {
	id, err := u.storage.Create(user)
	if err != nil {
		u.logger.Error("user: User.Create error", zap.Error(err))
	}

	return id, err
}

func (u *User) Exists(ctx context.Context, id int) error {
	_, err := u.storage.GetByID(id)

	return err
}
