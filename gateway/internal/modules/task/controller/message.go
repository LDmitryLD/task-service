package controller

type CreateTaskRequest struct {
	TaskName    string
	Description string
	UserID      int
}
