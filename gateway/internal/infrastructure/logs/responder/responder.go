package responder

import (
	"encoding/json"
	"net/http"
	"projects/LDmitryLD/task-service/gateway/internal/models"

	"go.uber.org/zap"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
	ErrorNotFound(w http.ResponseWriter, err error)
}

type Respond struct {
	log *zap.Logger
}

func NewResponder(logger *zap.Logger) Responder {
	return &Respond{log: logger}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
		r.ErrorInternal(w, err)
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.log.Error("http response internal error", zap.Error(err))

	w.Header().Set("Content-Type", "application/json")

	message := extractGrpcErrorMessage(err)
	resp := models.ApiResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}

	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		r.log.Error("error with encode internal error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))

	w.Header().Set("Content-Type", "application/json")

	message := extractGrpcErrorMessage(err)
	resp := models.ApiResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}

	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		r.log.Error("error with encode bad request error", zap.Error(err))
	}
}

func (r *Respond) ErrorNotFound(w http.ResponseWriter, err error) {
	r.log.Info("http response Not Found")

	w.Header().Set("Content-Type", "application/json")

	message := extractGrpcErrorMessage(err)
	resp := models.ApiResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}

	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		r.log.Error("error with encode not found error", zap.Error(err))
	}
}
