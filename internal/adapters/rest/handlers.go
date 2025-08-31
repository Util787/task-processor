package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Util787/task-processor/internal/common"
	"github.com/Util787/task-processor/internal/models"
)

type TaskUsecase interface {
	EnqueueTask(ctx context.Context, task models.Task) error
	GetTaskState(ctx context.Context, taskID string) (models.TaskState, error)
}

type Handler struct {
	log         *slog.Logger
	taskUsecase TaskUsecase
}

func (h *Handler) getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) enqueueTask(w http.ResponseWriter, r *http.Request) {
	log := logReqId(r.Context(), h.log)

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		newErrorResponse(w, log, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.taskUsecase.EnqueueTask(r.Context(), task); err != nil {
		if errors.Is(err, models.ErrValidation) {
			newErrorResponse(w, log, http.StatusBadRequest, "invalid task", err)
			return
		}
		newErrorResponse(w, log, http.StatusInternalServerError, "failed to enqueue task", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) getTaskState(w http.ResponseWriter, r *http.Request) {
	log := logReqId(r.Context(), h.log)

	taskID := r.URL.Query().Get("task_id")

	state, err := h.taskUsecase.GetTaskState(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, models.ErrTaskNotFound) {
			newErrorResponse(w, log, http.StatusNotFound, "task not found", err)
			return
		}
		newErrorResponse(w, log, http.StatusInternalServerError, "failed to get task state", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	type resp struct {
		TaskState models.TaskState `json:"task_state"`
	}

	if err := json.NewEncoder(w).Encode(resp{TaskState: state}); err != nil {
		newErrorResponse(w, log, http.StatusInternalServerError, "failed to encode response", err)
	}
}

// should be used in the start of every handler
func logReqId(ctx context.Context, log *slog.Logger) *slog.Logger {
	reqId := ctx.Value(common.ContextKey("request_id"))
	return log.With(slog.Any("request_id", reqId))
}
