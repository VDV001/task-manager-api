package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var _ domain.TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *domain.Task) error {
	const q = `
		INSERT INTO tasks (id, title, description, status, deadline, created_at, updated_at, author_id)
		VALUES (:id, :title, :description, :status, :deadline, :created_at, :updated_at, :author_id)`

	if _, err := r.db.NamedExecContext(ctx, q, task); err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	const q = `
		SELECT id, title, description, status, deadline, created_at, updated_at, deleted_at, author_id
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL`

	if err := r.db.GetContext(ctx, &task, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return &task, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {
	const q = `
		UPDATE tasks
		SET title = :title, description = :description, status = :status,
		    deadline = :deadline, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	res, err := r.db.NamedExecContext(ctx, q, task)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE tasks SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("soft delete task: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) List(ctx context.Context, authorID uuid.UUID, f domain.TaskFilter) ([]domain.Task, int, error) {
	var (
		conditions []string
		args       []any
		argIdx     int
	)

	nextArg := func(val any) string {
		argIdx++
		args = append(args, val)
		return fmt.Sprintf("$%d", argIdx)
	}

	conditions = append(conditions, "author_id = "+nextArg(authorID))
	conditions = append(conditions, "deleted_at IS NULL")

	if f.Status != nil {
		conditions = append(conditions, "status = "+nextArg(string(*f.Status)))
	}
	if f.Search != nil && *f.Search != "" {
		conditions = append(conditions, "title ILIKE "+nextArg("%"+*f.Search+"%"))
	}
	if f.Overdue != nil && *f.Overdue {
		conditions = append(conditions, "deadline < now()")
		conditions = append(conditions, "status != 'done'")
	}
	if f.DeadlineBefore != nil {
		conditions = append(conditions, "deadline <= "+nextArg(*f.DeadlineBefore))
	}
	if f.DeadlineAfter != nil {
		conditions = append(conditions, "deadline >= "+nextArg(*f.DeadlineAfter))
	}
	if f.CreatedAfter != nil {
		conditions = append(conditions, "created_at >= "+nextArg(*f.CreatedAfter))
	}
	if f.CreatedBefore != nil {
		conditions = append(conditions, "created_at <= "+nextArg(*f.CreatedBefore))
	}

	where := strings.Join(conditions, " AND ")

	// Подсчёт общего количества.
	countQ := "SELECT COUNT(*) FROM tasks WHERE " + where
	var total int
	if err := r.db.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	// Whitelist сортировки (SQL injection prevention).
	sortColumn, ok := map[string]string{
		"created_at": "created_at",
		"deadline":   "deadline",
		"status":     "status",
		"title":      "title",
	}[f.SortBy]
	if !ok {
		sortColumn = "created_at"
	}

	order := "DESC"
	if f.Order == "asc" {
		order = "ASC"
	}

	offset := (f.Page - 1) * f.Limit

	listQ := fmt.Sprintf(`
		SELECT id, title, description, status, deadline, created_at, updated_at, deleted_at, author_id
		FROM tasks
		WHERE %s
		ORDER BY %s %s
		LIMIT %s OFFSET %s`,
		where, sortColumn, order, nextArg(f.Limit), nextArg(offset),
	)

	var tasks []domain.Task
	if err := r.db.SelectContext(ctx, &tasks, listQ, args...); err != nil {
		return nil, 0, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, total, nil
}

func (r *TaskRepo) GetStats(ctx context.Context, authorID uuid.UUID) (*domain.TaskStats, error) {
	const q = `
		SELECT
			COUNT(*)                                         AS total,
			COUNT(*) FILTER (WHERE status = 'new')           AS new,
			COUNT(*) FILTER (WHERE status = 'in_progress')   AS in_progress,
			COUNT(*) FILTER (WHERE status = 'done')          AS done,
			COUNT(*) FILTER (WHERE deadline < now() AND status != 'done') AS overdue
		FROM tasks
		WHERE author_id = $1 AND deleted_at IS NULL`

	var row struct {
		Total      int `db:"total"`
		New        int `db:"new"`
		InProgress int `db:"in_progress"`
		Done       int `db:"done"`
		Overdue    int `db:"overdue"`
	}

	if err := r.db.GetContext(ctx, &row, q, authorID); err != nil {
		return nil, fmt.Errorf("get task stats: %w", err)
	}

	return &domain.TaskStats{
		Total: row.Total,
		ByStatus: map[string]int{
			"new":         row.New,
			"in_progress": row.InProgress,
			"done":        row.Done,
		},
		Overdue: row.Overdue,
	}, nil
}
