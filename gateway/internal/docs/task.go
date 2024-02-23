package docs

import (
	"projects/LDmitryLD/task-service/gateway/internal/models"
	taskcontroller "projects/LDmitryLD/task-service/gateway/internal/modules/task/controller"
)

// swagger:route POST /api/tasks task CreateTaskRequest
// Добавить задание.
// responses:
//  200: CreateTaskResponse

// swagger:parameters CreateTaskRequest
type CreateTaskRequest struct {
	// in:body
	// required: true
	Body taskcontroller.CreateTaskRequest
}

// swagger:response CreateTaskResponse
type CreateTaskResponse struct {
	// in:body
	Body models.ApiResponse
}

// swagger:route GET /api/tasks/{userId} task ListRequest
// Получить список заданий пользователя по его ID.
// responses:
//  200: ListResponse

// swagger:parameters ListRequest
type ListRequest struct {
	// ID пользователя.
	//
	// in:path
	// required: true
	UserID int `json:"userId"`
}

// swagger:response ListResponse
type ListResponse struct {
	// in:body
	Body []models.Task
}

// swagger:route DELETE /api/users/{userId}/tasks/{taskId} task DeleteRequest
// Удалить задание.
// responses:
//  200: DeleteResponse

// swagger:parameters DeleteRequest
type DeleteRequest struct {
	// ID пользователя.
	//
	// in:path
	// required: true
	UserID int `json:"userId"`
	// ID задания.
	//
	// in:path
	// required: true
	TaskID int `json:"taskId"`
}

// swagger:response DeleteResponse
type DeleteResponse struct {
	// in:body
	Body models.ApiResponse
}
