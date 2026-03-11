package handler

import (
	"time"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/google/uuid"
)

// --- Auth DTOs ---

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// --- Task DTOs ---

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=10000"`
	Deadline    *time.Time `json:"deadline,omitzero"`
}

type UpdateTaskRequest struct {
	Title       *string            `json:"title,omitzero" validate:"omitempty,min=1,max=255"`
	Description *string            `json:"description,omitzero" validate:"omitempty,max=10000"`
	Status      *domain.TaskStatus `json:"status,omitzero" validate:"omitempty,oneof=new in_progress done"`
	Deadline    *time.Time         `json:"deadline,omitzero"`
}

type TaskResponse struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status"`
	Deadline    *time.Time        `json:"deadline,omitzero"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	AuthorID    uuid.UUID         `json:"author_id"`
}

func TaskToResponse(t *domain.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		AuthorID:    t.AuthorID,
	}
}

func TasksToResponse(tasks []domain.Task) []TaskResponse {
	result := make([]TaskResponse, len(tasks))
	for i := range tasks {
		result[i] = TaskToResponse(&tasks[i])
	}
	return result
}

type StatsResponse struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	Overdue  int            `json:"overdue"`
}
