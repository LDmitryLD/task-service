package models

type TaskDTO struct {
	TaskID      int
	TaskName    string
	Description string
	UserID      int
}

type Task struct {
	TaskID      int
	TaskName    string
	Description string
}
