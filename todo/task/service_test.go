package task

import (
	"context"
	"errors"
	"fmt"
	"mkmgo-todo/todo/pagination"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	Mock task/repository.go
*/

type MockTaskRepository struct {
	WriteTaskFunc   func(ctx context.Context, task *Task) error
	GetAllTasksFunc func(ctx context.Context, request GetAllTaskRequest) ([]Task, error)
	DeleteTaskFunc  func(ctx context.Context, id uint64) error
}

func (m *MockTaskRepository) WriteTask(ctx context.Context, task *Task) error {
	if m.WriteTaskFunc != nil {
		return m.WriteTaskFunc(ctx, task)
	}
	return nil
}

func (m *MockTaskRepository) GetAllTasks(ctx context.Context, request GetAllTaskRequest) ([]Task, error) {
	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(ctx, request)
	}
	return []Task{}, nil
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id uint64) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}
	return nil
}

/*
	Unit test for task/service.go
*/

func TestWriteTask(t *testing.T) {
	mockRepo := &MockTaskRepository{
		WriteTaskFunc: func(ctx context.Context, task *Task) error {
			task.ID = 1
			task.UpdatedAt = time.Now()
			return nil
		},
	}
	service := NewTaskServiceImpl(mockRepo)

	req := &WriteTaskRequest{
		Title:       "New Task",
		Description: "Task Description",
	}
	resp, err := service.WriteTask(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Description, resp.Description)
}

func TestWriteTaskWhenFailAtRepoWriteTask(t *testing.T) {
	mockRepo := &MockTaskRepository{
		WriteTaskFunc: func(ctx context.Context, task *Task) error {
			return fmt.Errorf("write task error")
		},
	}

	service := NewTaskServiceImpl(mockRepo)

	req := &WriteTaskRequest{
		Title:       "New Task",
		Description: "Task Description",
	}

	resp, err := service.WriteTask(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "write task error", err.Error())
}

func TestGetAllTasks(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetAllTasksFunc: func(ctx context.Context, request GetAllTaskRequest) ([]Task, error) {
			return []Task{
				{ID: 1, Title: "Task 1", Description: "Description 1", UpdatedAt: time.Now()},
				{ID: 2, Title: "Task 2", Description: "Description 2", UpdatedAt: time.Now()},
			}, nil
		},
	}
	service := NewTaskServiceImpl(mockRepo)

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	resp, err := service.GetAllTasks(context.Background(), request)

	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "Task 1", resp[0].Title)
	assert.Equal(t, "Task 2", resp[1].Title)
}

func TestGetAllTaskFailAtRepoGetAllTask(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetAllTasksFunc: func(ctx context.Context, request GetAllTaskRequest) ([]Task, error) {
			return nil, fmt.Errorf("get all task error")
		},
	}

	service := NewTaskServiceImpl(mockRepo)

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	resp, err := service.GetAllTasks(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "get all task error", err.Error())
}

func TestDeleteTask(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteTaskFunc: func(ctx context.Context, id uint64) error {
			if id == 1 {
				return nil
			}
			return errors.New("task not found")
		},
	}
	service := NewTaskServiceImpl(mockRepo)

	err := service.DeleteTask(context.Background(), 1)
	assert.NoError(t, err)

	err = service.DeleteTask(context.Background(), 2)
	assert.Error(t, err)
}
