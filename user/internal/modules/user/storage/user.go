package storage

import (
	"projects/LDmitryLD/task-service/user/internal/db/adapter"
	"projects/LDmitryLD/task-service/user/internal/models"
)

type UserStorager interface {
	GetByID(id int) (models.User, error)
	Create(user models.User) (int, error)
}

type UserStorage struct {
	adapter adapter.SQLAdapterer
}

func NewUserStorage(sqlAdapter adapter.SQLAdapterer) UserStorager {
	return &UserStorage{
		adapter: sqlAdapter,
	}
}

func (u *UserStorage) GetByID(id int) (models.User, error) {
	return u.adapter.GetUserByID(id)
}

func (u *UserStorage) Create(user models.User) (int, error) {
	return u.adapter.CreateUser(user)
}
