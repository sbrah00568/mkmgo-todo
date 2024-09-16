package task

import (
	"context"
	"mkmgo-todo/todo/pagination"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSaveTaskMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	task := &Task{Title: "Mocked Task", Description: "Mocked Desc"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "task" ("title","description","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(task.Title, task.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.SaveTask(context.Background(), task)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), task.ID)
}

func TestSaveTaskMockWhenError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	task := &Task{Title: "Mocked Task", Description: "Mocked Desc"}

	mock.ExpectBegin()
	// WRONG QUERY
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("title","description","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(task.Title, task.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.SaveTask(context.Background(), task)

	assert.Error(t, err)
}

func TestGetAllTasksMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	expectedTasks := []Task{
		{ID: 1, Title: "Task 1", Description: "Description 1"},
		{ID: 2, Title: "Task 2", Description: "Description 2"},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "task"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Task 1", "Description 1", time.Now(), time.Now(), nil).
			AddRow(2, "Task 2", "Description 2", time.Now(), time.Now(), nil))

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	gotTasks, err := repo.GetAllTasks(context.Background(), request)

	assert.NoError(t, err)

	for i, task := range gotTasks {
		expectedTask := expectedTasks[i]
		assert.Equal(t, expectedTask.ID, task.ID)
		assert.Equal(t, expectedTask.Title, task.Title)
		assert.Equal(t, expectedTask.Description, task.Description)
	}
}

func TestGetAllTasksMockWhenError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	// WRONG QUERY
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Tasks 2", "Description X", time.Now(), time.Now(), nil))

	pagination := pagination.PaginationRequest{
		Page:     1,
		PageSize: 10,
		Order:    "title",
		SortBy:   "asc",
	}
	request := GetAllTaskRequest{
		PaginationRequest: &pagination,
	}
	_, err = repo.GetAllTasks(context.Background(), request)

	assert.Error(t, err)
}

func TestDeleteTaskMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "task" SET "deleted_at"=$1 WHERE "task"."id" = $2 AND "task"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteTask(context.Background(), 1)

	assert.NoError(t, err)
}

func TestDeleteTaskMockWhenError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepositoryImpl(gormDB)

	mock.ExpectBegin()
	// WRONG QUERY
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "deleted_at"=$1 WHERE "id" = $2 AND "deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteTask(context.Background(), 1)

	assert.Error(t, err)
}
