package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"mkmgo-todo/todo/pagination"
	"mkmgo-todo/todo/task"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskService interface {
	WriteTask(ctx context.Context, request *task.WriteTaskRequest) (*task.GetTaskResponse, error)
	GetAllTasks(ctx context.Context, request task.GetAllTaskRequest) ([]task.GetTaskResponse, error)
	DeleteTask(ctx context.Context, id uint64) error
}

type TaskHandler struct {
	taskSvc TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{taskSvc: service}
}

func (h *TaskHandler) WriteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req task.WriteTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}
	res, err := h.taskSvc.WriteTask(r.Context(), &req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, res)
}

func (h *TaskHandler) GetAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	pagination := pagination.NewPaginationRequest(r)
	request := task.GetAllTaskRequest{PaginationRequest: pagination}
	res, err := h.taskSvc.GetAllTasks(r.Context(), request)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, res)
}

func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.taskSvc.DeleteTask(r.Context(), id); err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, fmt.Sprintf("Task %d deleted", id))
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
