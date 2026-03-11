package usecase_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTaskUseCase(taskRepo domain.TaskRepository) *usecase.TaskUseCase {
	return usecase.NewTaskUseCase(taskRepo, slog.Default())
}

func TestTaskUseCase_Create(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	authorID := uuid.New()

	repo := &mockTaskRepo{
		createFn: func(_ context.Context, task *domain.Task) error {
			assert.Equal(t, "Deploy v2", task.Title)
			assert.Equal(t, domain.TaskStatusNew, task.Status)
			assert.Equal(t, authorID, task.AuthorID)
			assert.NotEmpty(t, task.ID)
			return nil
		},
	}

	uc := newTaskUseCase(repo)
	deadline := time.Now().Add(24 * time.Hour)
	task, err := uc.Create(ctx, authorID, usecase.CreateTaskInput{
		Title:       "Deploy v2",
		Description: "Deploy to production",
		Deadline:    &deadline,
	})

	require.NoError(t, err)
	assert.Equal(t, "Deploy v2", task.Title)
	assert.Equal(t, domain.TaskStatusNew, task.Status)
	assert.Equal(t, authorID, task.AuthorID)
}

func TestTaskUseCase_GetByID(t *testing.T) {
	t.Parallel()

	authorID := uuid.New()
	otherUserID := uuid.New()
	taskID := uuid.New()

	existingTask := &domain.Task{
		ID:       taskID,
		Title:    "Test",
		AuthorID: authorID,
	}

	tests := []struct {
		name     string
		callerID uuid.UUID
		wantErr  error
	}{
		{
			name:     "owner can access",
			callerID: authorID,
			wantErr:  nil,
		},
		{
			name:     "non-owner forbidden",
			callerID: otherUserID,
			wantErr:  domain.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repo := &mockTaskRepo{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
					return existingTask, nil
				},
			}

			uc := newTaskUseCase(repo)
			task, err := uc.GetByID(ctx, tt.callerID, taskID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, taskID, task.ID)
			}
		})
	}
}

func TestTaskUseCase_Update(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	authorID := uuid.New()
	taskID := uuid.New()

	repo := &mockTaskRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
			return &domain.Task{
				ID:       taskID,
				Title:    "Old title",
				Status:   domain.TaskStatusNew,
				AuthorID: authorID,
			}, nil
		},
		updateFn: func(_ context.Context, task *domain.Task) error {
			assert.Equal(t, "New title", task.Title)
			assert.Equal(t, domain.TaskStatusInProgress, task.Status)
			return nil
		},
	}

	uc := newTaskUseCase(repo)
	newTitle := "New title"
	newStatus := domain.TaskStatusInProgress

	task, err := uc.Update(ctx, authorID, taskID, usecase.UpdateTaskInput{
		Title:  &newTitle,
		Status: &newStatus,
	})

	require.NoError(t, err)
	assert.Equal(t, "New title", task.Title)
	assert.Equal(t, domain.TaskStatusInProgress, task.Status)
}

func TestTaskUseCase_Update_Forbidden(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	otherUserID := uuid.New()
	taskID := uuid.New()

	repo := &mockTaskRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
			return &domain.Task{
				ID:       taskID,
				AuthorID: uuid.New(), // другой пользователь
			}, nil
		},
	}

	uc := newTaskUseCase(repo)
	newTitle := "Hacked"
	_, err := uc.Update(ctx, otherUserID, taskID, usecase.UpdateTaskInput{Title: &newTitle})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTaskUseCase_Delete(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	authorID := uuid.New()
	taskID := uuid.New()

	deleted := false
	repo := &mockTaskRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
			return &domain.Task{ID: taskID, AuthorID: authorID}, nil
		},
		softDeleteFn: func(_ context.Context, id uuid.UUID) error {
			assert.Equal(t, taskID, id)
			deleted = true
			return nil
		},
	}

	uc := newTaskUseCase(repo)
	err := uc.Delete(ctx, authorID, taskID)

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestTaskUseCase_List_DefaultPagination(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	authorID := uuid.New()

	repo := &mockTaskRepo{
		listFn: func(_ context.Context, _ uuid.UUID, filter *domain.TaskFilter) ([]domain.Task, int, error) {
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 20, filter.Limit)
			assert.Equal(t, "created_at", filter.SortBy)
			assert.Equal(t, "desc", filter.Order)
			return []domain.Task{}, 0, nil
		},
	}

	uc := newTaskUseCase(repo)
	tasks, total, err := uc.List(ctx, authorID, &domain.TaskFilter{})

	require.NoError(t, err)
	assert.Empty(t, tasks)
	assert.Equal(t, 0, total)
}

func TestTaskUseCase_List_SanitizesInput(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	authorID := uuid.New()

	repo := &mockTaskRepo{
		listFn: func(_ context.Context, _ uuid.UUID, filter *domain.TaskFilter) ([]domain.Task, int, error) {
			assert.Equal(t, "created_at", filter.SortBy, "invalid sort_by should fallback to created_at")
			assert.Equal(t, "desc", filter.Order, "invalid order should fallback to desc")
			assert.Equal(t, 20, filter.Limit, "limit > 100 should fallback to 20")
			return nil, 0, nil
		},
	}

	uc := newTaskUseCase(repo)
	_, _, err := uc.List(ctx, authorID, &domain.TaskFilter{
		SortBy: "DROP TABLE tasks;--",
		Order:  "sideways",
		Limit:  999,
	})
	require.NoError(t, err)
}

func TestTaskUseCase_Stats(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	authorID := uuid.New()

	expectedStats := &domain.TaskStats{
		Total: 10,
		ByStatus: map[string]int{
			"new":         3,
			"in_progress": 4,
			"done":        3,
		},
		Overdue: 2,
	}

	repo := &mockTaskRepo{
		getStatsFn: func(_ context.Context, id uuid.UUID) (*domain.TaskStats, error) {
			assert.Equal(t, authorID, id)
			return expectedStats, nil
		},
	}

	uc := newTaskUseCase(repo)
	stats, err := uc.Stats(ctx, authorID)

	require.NoError(t, err)
	assert.Equal(t, 10, stats.Total)
	assert.Equal(t, 2, stats.Overdue)
	assert.Equal(t, 4, stats.ByStatus["in_progress"])
}
