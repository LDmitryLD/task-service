package docs

import (
	"projects/LDmitryLD/task-service/gateway/internal/models"
	usercontroller "projects/LDmitryLD/task-service/gateway/internal/modules/user/controller"
)

// swagger:route POST /api/users user CreateRequest
// Добавить пользователя.
// responses:
//  200: CreateResponse

// swagger:parameters CreateRequest
type CreateRequest struct {
	// in:body
	// required: true
	Body usercontroller.CreateRequest
}

// swagger:response CreateResponse
type CreateResponse struct {
	// in:body
	Body models.ApiResponse
}

// swagger:route GET /api/users/{id} user ProfileRequest
// Получить пользователя по ID.
// responses:
//  200: ProfileResponse

// swagger:parameters ProfileRequest
type ProfileRequest struct {
	// ID пользователя.
	//
	// in:path
	// required: true
	ID int `json:"id"`
}

// swagger:response ProfileResponse
type ProfielResponse struct {
	// in:body
	Body models.User
}
