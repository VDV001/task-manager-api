package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "new"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusNew, TaskStatusInProgress, TaskStatusDone:
		return true
	}
	return false
}

type Task struct {
	ID          uuid.UUID  `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      TaskStatus `db:"status"`
	Deadline    *time.Time `db:"deadline"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	AuthorID    uuid.UUID  `db:"author_id"`
}

func (t *Task) IsOverdue() bool {
	if t.Deadline == nil || t.Status == TaskStatusDone {
		return false
	}
	return time.Now().After(*t.Deadline)
}

// TaskFilter — параметры фильтрации и пагинации списка задач.
type TaskFilter struct {
	Status         *TaskStatus
	Search         *string
	Overdue        *bool
	DeadlineBefore *time.Time
	DeadlineAfter  *time.Time
	CreatedAfter   *time.Time
	CreatedBefore  *time.Time
	SortBy         string // created_at, deadline, status, title
	Order          string // asc, desc
	Page           int
	Limit          int
}

// TaskStats — агрегированная статистика по задачам пользователя.
type TaskStats struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	Overdue  int            `json:"overdue"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, authorID uuid.UUID, filter *TaskFilter) ([]Task, int, error)
	GetStats(ctx context.Context, authorID uuid.UUID) (*TaskStats, error)
}
