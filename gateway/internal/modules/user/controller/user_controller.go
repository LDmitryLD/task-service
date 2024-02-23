package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/component"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/logs/responder"
	"projects/LDmitryLD/task-service/gateway/internal/models"
	"projects/LDmitryLD/task-service/gateway/internal/modules/user/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Userer interface {
	Profile(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	user service.Userer
	responder.Responder
}

func NewUserController(service service.Userer, components *component.Components) Userer {
	return &UserController{
		user:      service,
		Responder: components.Responder,
	}
}

func (u *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	idRaw := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	user, err := u.user.Profile(r.Context(), id)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			u.ErrorNotFound(w, err)
			return
		default:
			u.ErrorInternal(w, err)
			return
		}
	}

	u.OutputJSON(w, user)

}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	id, err := u.user.Create(r.Context(), models.User{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email})
	if err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	resp := models.ApiResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("пользователь с id %d добавлен", id),
	}

	u.OutputJSON(w, resp)
}
