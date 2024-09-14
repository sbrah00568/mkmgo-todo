package task

import "context"

type TaskRepository interface {
	WriteTask(ctx context.Context, task *Task) error
	GetAllTasks(ctx context.Context) ([]Task, error)
	DeleteTask(ctx context.Context, id uint64) error
}

type TaskServiceImpl struct {
	repo TaskRepository
}

func NewTaskServiceImpl(repo TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{repo: repo}
}

func (svc *TaskServiceImpl) WriteTask(ctx context.Context, req *WriteTaskRequest) (*GetTaskResponse, error) {
	task := Task{
		Title:       req.Title,
		Description: req.Description,
	}
	err := svc.repo.WriteTask(ctx, &task)
	if err != nil {
		return nil, err
	}
	return &GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (svc *TaskServiceImpl) GetAllTasks(ctx context.Context) ([]GetTaskResponse, error) {
	tasks, err := svc.repo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	var responses []GetTaskResponse
	for _, task := range tasks {
		responses = append(responses, GetTaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			UpdatedAt:   task.UpdatedAt,
		})
	}
	return responses, nil
}

func (svc *TaskServiceImpl) DeleteTask(ctx context.Context, id uint64) error {
	return svc.repo.DeleteTask(ctx, id)
}
