package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/daniilgit/task-manager-api/pkg/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TaskHandler struct {
	tasks    *usecase.TaskUseCase
	validate *validator.Validate
}

func NewTaskHandler(tasks *usecase.TaskUseCase, validate *validator.Validate) *TaskHandler {
	return &TaskHandler{tasks: tasks, validate: validate}
}

// Create godoc
// @Summary      Создать задачу
// @Description  Создаёт новую задачу для текущего пользователя.
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body CreateTaskRequest true "Данные задачи"
// @Success      201 {object} httputil.Response{data=TaskResponse}
// @Failure      400 {object} httputil.ErrorResponse "Невалидный JSON"
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      422 {object} httputil.ErrorResponse "Ошибка валидации"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if errs := validateStruct(h.validate, req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}
	task, err := h.tasks.Create(r.Context(), authorID, usecase.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
	})
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create task")
		return
	}

	httputil.Success(w, http.StatusCreated, TaskToResponse(task))
}

// GetByID godoc
// @Summary      Получить задачу по ID
// @Description  Возвращает задачу по UUID. Доступна только автору.
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Task ID (UUID)"
// @Success      200 {object} httputil.Response{data=TaskResponse}
// @Failure      400 {object} httputil.ErrorResponse "Невалидный UUID"
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      403 {object} httputil.ErrorResponse "Доступ запрещён"
// @Failure      404 {object} httputil.ErrorResponse "Задача не найдена"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid task ID")
		return
	}

	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}
	task, err := h.tasks.GetByID(r.Context(), authorID, taskID)
	if err != nil {
		handleTaskError(w, err)
		return
	}

	httputil.Success(w, http.StatusOK, TaskToResponse(task))
}

// Update godoc
// @Summary      Обновить задачу
// @Description  Частичное обновление задачи (PATCH). Доступна только автору.
// @Tags         tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID (UUID)"
// @Param        body body UpdateTaskRequest true "Поля для обновления"
// @Success      200 {object} httputil.Response{data=TaskResponse}
// @Failure      400 {object} httputil.ErrorResponse "Невалидный JSON или UUID"
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      403 {object} httputil.ErrorResponse "Доступ запрещён"
// @Failure      404 {object} httputil.ErrorResponse "Задача не найдена"
// @Failure      422 {object} httputil.ErrorResponse "Ошибка валидации"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks/{id} [patch]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if errs := validateStruct(h.validate, req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}
	task, err := h.tasks.Update(r.Context(), authorID, taskID, usecase.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Deadline:    req.Deadline,
	})
	if err != nil {
		handleTaskError(w, err)
		return
	}

	httputil.Success(w, http.StatusOK, TaskToResponse(task))
}

// Delete godoc
// @Summary      Удалить задачу (soft delete)
// @Description  Мягкое удаление задачи. Доступна только автору.
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Task ID (UUID)"
// @Success      204 "Задача удалена"
// @Failure      400 {object} httputil.ErrorResponse "Невалидный UUID"
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      403 {object} httputil.ErrorResponse "Доступ запрещён"
// @Failure      404 {object} httputil.ErrorResponse "Задача не найдена"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid task ID")
		return
	}

	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}
	if err := h.tasks.Delete(r.Context(), authorID, taskID); err != nil {
		handleTaskError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary      Список задач с фильтрацией и пагинацией
// @Description  Возвращает список задач текущего пользователя с фильтрацией, сортировкой и пагинацией.
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Param        status query string false "Фильтр по статусу" Enums(new, in_progress, done)
// @Param        search query string false "Поиск по заголовку"
// @Param        overdue query bool false "Только просроченные"
// @Param        deadline_before query string false "Дедлайн до (YYYY-MM-DD)"
// @Param        deadline_after query string false "Дедлайн после (YYYY-MM-DD)"
// @Param        created_after query string false "Создано после (YYYY-MM-DD)"
// @Param        created_before query string false "Создано до (YYYY-MM-DD)"
// @Param        sort_by query string false "Поле сортировки" Enums(created_at, deadline, status, title)
// @Param        order query string false "Направление" Enums(asc, desc)
// @Param        page query int false "Страница" default(1)
// @Param        limit query int false "Лимит" default(20)
// @Success      200 {object} httputil.Response{data=[]TaskResponse,meta=httputil.PaginationMeta}
// @Failure      400 {object} httputil.ErrorResponse "Невалидный параметр фильтрации"
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := domain.TaskFilter{
		SortBy: q.Get("sort_by"),
		Order:  q.Get("order"),
		Page:   intParam(q.Get("page"), 1),
		Limit:  intParam(q.Get("limit"), 20),
	}

	if s := q.Get("status"); s != "" {
		status := domain.TaskStatus(s)
		if !status.IsValid() {
			httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid status filter, must be one of: new, in_progress, done")
			return
		}
		filter.Status = &status
	}
	if s := q.Get("search"); s != "" {
		filter.Search = &s
	}
	if s := q.Get("overdue"); s == "true" {
		overdue := true
		filter.Overdue = &overdue
	}
	if t := parseTime(q.Get("deadline_before")); t != nil {
		filter.DeadlineBefore = t
	}
	if t := parseTime(q.Get("deadline_after")); t != nil {
		filter.DeadlineAfter = t
	}
	if t := parseTime(q.Get("created_after")); t != nil {
		filter.CreatedAfter = t
	}
	if t := parseTime(q.Get("created_before")); t != nil {
		filter.CreatedBefore = t
	}

	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}
	tasks, total, err := h.tasks.List(r.Context(), authorID, &filter)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list tasks")
		return
	}

	httputil.SuccessWithMeta(w, http.StatusOK, TasksToResponse(tasks), httputil.PaginationMeta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	})
}

// Stats godoc
// @Summary      Статистика по задачам
// @Description  Возвращает агрегированную статистику: общее количество, по статусам, просроченные.
// @Tags         tasks
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} httputil.Response{data=StatsResponse}
// @Failure      401 {object} httputil.ErrorResponse "Не авторизован"
// @Failure      500 {object} httputil.ErrorResponse "Внутренняя ошибка"
// @Router       /tasks/stats [get]
func (h *TaskHandler) Stats(w http.ResponseWriter, r *http.Request) {
	authorID, err := UserIDFromContext(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing user identity")
		return
	}

	stats, err := h.tasks.Stats(r.Context(), authorID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get stats")
		return
	}

	httputil.Success(w, http.StatusOK, StatsResponse{
		Total:    stats.Total,
		ByStatus: stats.ByStatus,
		Overdue:  stats.Overdue,
	})
}

func handleTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		httputil.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
	case errors.Is(err, domain.ErrForbidden):
		httputil.Error(w, http.StatusForbidden, "FORBIDDEN", "you can only access your own tasks")
	default:
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func intParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	// Поддержка формата даты (YYYY-MM-DD) и полного ISO 8601.
	for _, layout := range []string{time.DateOnly, time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
