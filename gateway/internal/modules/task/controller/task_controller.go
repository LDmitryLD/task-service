package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/component"
	"projects/LDmitryLD/task-service/gateway/internal/infrastructure/logs/responder"
	"projects/LDmitryLD/task-service/gateway/internal/models"
	"projects/LDmitryLD/task-service/gateway/internal/modules/task/service"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Tasker interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type TaskController struct {
	task service.Tasker
	responder.Responder
}

func NewTaskController(task service.Tasker, components *component.Components) Tasker {
	return &TaskController{
		task:      task,
		Responder: components.Responder,
	}
}

func (t *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.ErrorBadRequest(w, err)
		return
	}

	id, err := t.task.Create(r.Context(), models.TaskDTO{TaskName: req.TaskName, Description: req.Description, UserID: req.UserID})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			t.ErrorNotFound(w, err)
			return
		default:
			t.ErrorInternal(w, err)
			return
		}
	}

	resp := models.ApiResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("задание с id %d добавлено", id),
	}

	t.OutputJSON(w, resp)
}

func (t *TaskController) List(w http.ResponseWriter, r *http.Request) {
	idRaw := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		t.ErrorBadRequest(w, err)
		return
	}

	tasks, err := t.task.List(r.Context(), id)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			t.ErrorNotFound(w, err)
			return
		default:
			t.ErrorInternal(w, err)
			return
		}
	}

	t.OutputJSON(w, tasks)
}

func (t *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	taskIdRaw := chi.URLParam(r, "taskId")
	taskId, err := strconv.Atoi(taskIdRaw)
	if err != nil {
		t.ErrorBadRequest(w, err)
		return
	}

	userIdRaw := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdRaw)
	if err != nil {
		t.ErrorBadRequest(w, err)
		return
	}

	_, err = t.task.Delete(r.Context(), userId, taskId)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			t.ErrorNotFound(w, err)
			return
		default:
			t.ErrorInternal(w, err)
			return
		}
	}

	resp := models.ApiResponse{
		Code:    200,
		Message: "задание удалено",
	}
	t.OutputJSON(w, resp)
}

// func extractGrpcErrorMessage(err error) string {
// 	if st, ok := status.FromError(err); ok {
// 		return st.Message()
// 	}

// 	return err.Error()
// }
