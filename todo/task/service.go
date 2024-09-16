package task

import (
	"context"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, task *Task) error
	GetAllTasks(ctx context.Context, request GetAllTaskRequest) ([]Task, error)
	DeleteTask(ctx context.Context, id uint64) error
}

type TaskServiceImpl struct {
	repo TaskRepository
}

func NewTaskServiceImpl(repo TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{repo: repo}
}

func (svc *TaskServiceImpl) SaveTask(ctx context.Context, request *WriteTaskRequest) (*GetTaskResponse, error) {
	task := Task{
		Title:       request.Title,
		Description: request.Description,
	}
	if request.ID != 0 {
		task.ID = request.ID
	}
	if err := svc.repo.SaveTask(ctx, &task); err != nil {
		return nil, err
	}
	return &GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		UpdatedAt:   task.FormattedUpdatedAt(),
	}, nil
}

func (svc *TaskServiceImpl) GetAllTasks(ctx context.Context, request GetAllTaskRequest) ([]GetTaskResponse, error) {
	tasks, err := svc.repo.GetAllTasks(ctx, request)
	if err != nil {
		return nil, err
	}
	responses := make([]GetTaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = GetTaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			UpdatedAt:   task.FormattedUpdatedAt(),
		}
	}
	return responses, nil
}

func (svc *TaskServiceImpl) DeleteTask(ctx context.Context, id uint64) error {
	if err := svc.repo.DeleteTask(ctx, id); err != nil {
		return err
	}
	return nil
}
