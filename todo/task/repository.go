package task

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewTaskRepositoryImpl(db *gorm.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{DB: db}
}

func (r *TaskRepositoryImpl) SaveTask(ctx context.Context, task *Task) error {
	log := zerolog.Ctx(ctx).With().Str("method", "taskService.saveTask").Logger()
	if err := r.DB.WithContext(ctx).Save(task).Error; err != nil {
		log.Error().Err(err).Msg("failed to save task")
		return fmt.Errorf("failed to save task: %w", err)
	}
	log.Info().Msg("success to save task")
	return nil
}

func (r *TaskRepositoryImpl) GetAllTasks(ctx context.Context, request GetAllTaskRequest) ([]Task, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "taskService.GetAllTasks").Logger()
	var tasks []Task
	err := r.DB.WithContext(ctx).Model(&Task{}).
		Limit(request.PaginationRequest.PageSize).
		Offset(request.PaginationRequest.GetOffset()).
		Find(&tasks).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve tasks")
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	log.Info().Msg("success to retrieve tasks")
	return tasks, err
}

func (r *TaskRepositoryImpl) DeleteTask(ctx context.Context, id uint64) error {
	log := zerolog.Ctx(ctx).With().Str("method", "taskService.DeleteTask").Logger()
	if err := r.DB.WithContext(ctx).Delete(&Task{}, id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete task")
		return fmt.Errorf("failed to delete task: %w", err)
	}
	log.Info().Msg("success to delete task")
	return nil
}
