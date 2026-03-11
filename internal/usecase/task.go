package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/google/uuid"
)

type CreateTaskInput struct {
	Title       string
	Description string
	Deadline    *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *domain.TaskStatus
	Deadline    *time.Time
}

type TaskUseCase struct {
	tasks domain.TaskRepository
	log   *slog.Logger
}

func NewTaskUseCase(tasks domain.TaskRepository, log *slog.Logger) *TaskUseCase {
	return &TaskUseCase{
		tasks: tasks,
		log:   log,
	}
}

func (uc *TaskUseCase) Create(ctx context.Context, authorID uuid.UUID, input CreateTaskInput) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.TaskStatusNew,
		Deadline:    input.Deadline,
		CreatedAt:   now,
		UpdatedAt:   now,
		AuthorID:    authorID,
	}

	if err := uc.tasks.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	uc.log.Info("task created", slog.String("task_id", task.ID.String()))
	return task, nil
}

func (uc *TaskUseCase) GetByID(ctx context.Context, authorID, taskID uuid.UUID) (*domain.Task, error) {
	task, err := uc.tasks.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.AuthorID != authorID {
		return nil, domain.ErrForbidden
	}

	return task, nil
}

func (uc *TaskUseCase) Update(ctx context.Context, authorID, taskID uuid.UUID, input UpdateTaskInput) (*domain.Task, error) {
	task, err := uc.tasks.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.AuthorID != authorID {
		return nil, domain.ErrForbidden
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		if !input.Status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", *input.Status)
		}
		task.Status = *input.Status
	}
	if input.Deadline != nil {
		task.Deadline = input.Deadline
	}

	task.UpdatedAt = time.Now()

	if err := uc.tasks.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

func (uc *TaskUseCase) Delete(ctx context.Context, authorID, taskID uuid.UUID) error {
	task, err := uc.tasks.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if task.AuthorID != authorID {
		return domain.ErrForbidden
	}

	if err := uc.tasks.SoftDelete(ctx, taskID); err != nil {
		return fmt.Errorf("soft delete task: %w", err)
	}

	uc.log.Info("task deleted", slog.String("task_id", taskID.String()))
	return nil
}

func (uc *TaskUseCase) List(ctx context.Context, authorID uuid.UUID, filter *domain.TaskFilter) ([]domain.Task, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	allowedSortFields := map[string]bool{
		"created_at": true, "deadline": true, "status": true, "title": true,
	}
	if !allowedSortFields[filter.SortBy] {
		filter.SortBy = "created_at"
	}

	if filter.Order != "asc" && filter.Order != "desc" {
		filter.Order = "desc"
	}

	return uc.tasks.List(ctx, authorID, filter)
}

func (uc *TaskUseCase) Stats(ctx context.Context, authorID uuid.UUID) (*domain.TaskStats, error) {
	return uc.tasks.GetStats(ctx, authorID)
}
