package task

import (
	"mkmgo-todo/todo/pagination"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"not null"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"not null"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (Task) TableName() string {
	return "task"
}

type WriteTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetTaskResponse struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetAllTaskRequest struct {
	PaginationRequest *pagination.PaginationRequest
}
