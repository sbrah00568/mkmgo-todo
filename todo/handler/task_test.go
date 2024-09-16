package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mkmgo-todo/todo/pagination"
	"mkmgo-todo/todo/task"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

/*
	Mock task/service.go
*/

type MockTaskService struct {
	WriteTaskFunc   func(ctx context.Context, request *task.WriteTaskRequest) (*task.GetTaskResponse, error)
	GetAllTasksFunc func(ctx context.Context, request task.GetAllTaskRequest) ([]task.GetTaskResponse, error)
	DeleteTaskFunc  func(ctx context.Context, id uint64) error
}

func (m *MockTaskService) WriteTask(ctx context.Context, request *task.WriteTaskRequest) (*task.GetTaskResponse, error) {
	if m.WriteTaskFunc != nil {
		return m.WriteTaskFunc(ctx, request)
	}
	return nil, nil
}

func (m *MockTaskService) GetAllTasks(ctx context.Context, request task.GetAllTaskRequest) ([]task.GetTaskResponse, error) {
	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(ctx, request)
	}
	return nil, nil
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id uint64) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}
	return nil
}

/*
	Unit test for handler/task.go
*/

const (
	testTitle               = "Makima"
	testDescription         = "Makima super kawaii"
	testID                  = 1
	invalidWriteTaskRequest = `{"title":"Makima","description"}`
	validWriteTaskRequest   = `{"title":"Makima","description":"Makima super kawaii"}`
	validGetAllTaskRequest  = `{"page":1,"pageSize":10, "sortBy":"title", "orderBy":"asc"}`
	tasksUrl                = "/tasks"
)

func TestWriteTaskHandler(t *testing.T) {
	mockService := &MockTaskService{
		WriteTaskFunc: func(ctx context.Context, request *task.WriteTaskRequest) (*task.GetTaskResponse, error) {
			response := task.GetTaskResponse{
				ID:          testID,
				Title:       testTitle,
				Description: testDescription,
				UpdatedAt:   time.Now(),
			}
			return &response, nil
		},
	}

	handler := NewTaskHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, tasksUrl, bytes.NewBufferString(validWriteTaskRequest))
	w := httptest.NewRecorder()
	handler.WriteTaskHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, testTitle, respBody["title"])
	assert.Equal(t, testDescription, respBody["description"])
}

func TestWriteTaskHandlerWhenInvalidJSONRequest(t *testing.T) {
	mockService := &MockTaskService{}
	handler := NewTaskHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, tasksUrl, bytes.NewBufferString(invalidWriteTaskRequest))
	w := httptest.NewRecorder()
	handler.WriteTaskHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	var respBody map[string]string
	err := json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Error(t, err)
}

func TestWriteTaskHandlerWhenSvcWriteTaskFail(t *testing.T) {
	mockService := &MockTaskService{
		WriteTaskFunc: func(ctx context.Context, request *task.WriteTaskRequest) (*task.GetTaskResponse, error) {
			return nil, fmt.Errorf("write task error")
		},
	}

	handler := NewTaskHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, tasksUrl, bytes.NewBufferString(validWriteTaskRequest))
	w := httptest.NewRecorder()
	handler.WriteTaskHandler(w, r)

	req := task.WriteTaskRequest{
		Title:       testTitle,
		Description: testDescription,
	}
	_, err := handler.taskSvc.WriteTask(r.Context(), &req)
	assert.Error(t, err)
}

func TestGetAllTaskHandler(t *testing.T) {
	mockService := &MockTaskService{
		GetAllTasksFunc: func(ctx context.Context, request task.GetAllTaskRequest) ([]task.GetTaskResponse, error) {
			var responses []task.GetTaskResponse
			response := task.GetTaskResponse{
				ID:          testID,
				Title:       testTitle,
				Description: testDescription,
				UpdatedAt:   time.Now(),
			}
			responses = append(responses, response)
			return responses, nil
		},
	}

	expectedResponses := []task.GetTaskResponse{
		{ID: 1, Title: testTitle, Description: testDescription},
	}

	handler := NewTaskHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, tasksUrl, bytes.NewBufferString(validGetAllTaskRequest))
	w := httptest.NewRecorder()
	handler.GetAllTaskHandler(w, r)

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := task.GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	responses, err := handler.taskSvc.GetAllTasks(context.Background(), request)
	assert.NoError(t, err)

	for i, response := range responses {
		expectedResponse := expectedResponses[i]
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Title, response.Title)
		assert.Equal(t, expectedResponse.Description, response.Description)
	}

}

func TestGetAllTaskHandlerWhenSvcGetAllFail(t *testing.T) {
	mockService := &MockTaskService{
		GetAllTasksFunc: func(ctx context.Context, request task.GetAllTaskRequest) ([]task.GetTaskResponse, error) {
			return nil, fmt.Errorf("get all task error")
		},
	}

	handler := NewTaskHandler(mockService)

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := task.GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	r := httptest.NewRequest(http.MethodGet, tasksUrl, bytes.NewBufferString(validGetAllTaskRequest))
	w := httptest.NewRecorder()
	handler.GetAllTaskHandler(w, r)

	_, err := handler.taskSvc.GetAllTasks(context.Background(), request)
	assert.Error(t, err)
}

func TestDeleteTaskHandler(t *testing.T) {
	mockService := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, id uint64) error {
			return nil
		},
	}

	handler := NewTaskHandler(mockService)
	r := httptest.NewRequest(http.MethodDelete, tasksUrl+"/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler.DeleteTaskHandler(w, r)

	err := handler.taskSvc.DeleteTask(context.Background(), 1)
	assert.NoError(t, err)
}

func TestDeleteTaskHandlerInvalidID(t *testing.T) {
	mockService := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, id uint64) error {
			return fmt.Errorf("invalid id")
		},
	}

	handler := NewTaskHandler(mockService)
	r := httptest.NewRequest(http.MethodDelete, tasksUrl+"/1", nil)
	w := httptest.NewRecorder()
	handler.DeleteTaskHandler(w, r)

	err := handler.taskSvc.DeleteTask(context.Background(), 1)
	assert.Error(t, err)
}

func TestDeleteTaskHandlerWhenSvcFail(t *testing.T) {
	mockService := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, id uint64) error {
			return fmt.Errorf("delete task error")
		},
	}

	handler := NewTaskHandler(mockService)
	r := httptest.NewRequest(http.MethodDelete, tasksUrl+"/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler.DeleteTaskHandler(w, r)

	err := handler.taskSvc.DeleteTask(context.Background(), 1)
	assert.Error(t, err)
}
