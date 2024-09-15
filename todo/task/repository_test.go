package task

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWriteTaskMock(t *testing.T) {
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
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Returning ID as 1
	mock.ExpectCommit()

	err = repo.WriteTask(context.Background(), task)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), task.ID)
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

	gotTasks, err := repo.GetAllTasks(context.Background())

	assert.NoError(t, err)

	for i, task := range gotTasks {
		expectedTask := expectedTasks[i]
		assert.Equal(t, expectedTask.ID, task.ID)
		assert.Equal(t, expectedTask.Title, task.Title)
		assert.Equal(t, expectedTask.Description, task.Description)
	}
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
