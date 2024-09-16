package task

import (
	"context"

	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewTaskRepositoryImpl(db *gorm.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{DB: db}
}

func (r *TaskRepositoryImpl) WriteTask(ctx context.Context, task *Task) error {
	return r.DB.WithContext(ctx).Create(task).Error
}

func (r *TaskRepositoryImpl) GetAllTasks(ctx context.Context, request GetAllTaskRequest) ([]Task, error) {
	var tasks []Task
	err := r.DB.WithContext(ctx).Model(&Task{}).
		Limit(request.PaginationRequest.PageSize).
		Offset(request.PaginationRequest.GetOffset()).
		Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepositoryImpl) DeleteTask(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).Delete(&Task{}, id).Error
}
