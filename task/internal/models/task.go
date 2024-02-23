package models

type Task struct {
	TaskId      int    `json:"task_id" db:"id"`
	TaskName    string `db:"task_name"`
	Description string `db:"description"`
	UserId      int    `json:"user_id" db:"user_id"`
}
