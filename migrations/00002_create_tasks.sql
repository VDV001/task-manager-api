-- +goose Up
CREATE TYPE task_status AS ENUM ('new', 'in_progress', 'done');

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255)  NOT NULL,
    description TEXT          NOT NULL DEFAULT '',
    status      task_status   NOT NULL DEFAULT 'new',
    deadline    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    author_id   UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Partial index: активные задачи (soft delete).
CREATE INDEX idx_tasks_author_active ON tasks (author_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Фильтрация по статусу.
CREATE INDEX idx_tasks_status ON tasks (author_id, status)
    WHERE deleted_at IS NULL;

-- Фильтрация по дедлайну (просроченные задачи).
CREATE INDEX idx_tasks_deadline ON tasks (author_id, deadline)
    WHERE deleted_at IS NULL AND deadline IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS tasks;
DROP TYPE IF EXISTS task_status;
