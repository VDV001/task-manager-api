# Task Manager API

REST API для приложения "Менеджер задач" с аутентификацией, управлением задачами и статистикой.

## Архитектура

Проект построен на принципах **DDD (Domain-Driven Design)** и **Clean Architecture**:

```
cmd/api/              → Точка входа, сборка зависимостей (DI)
internal/
  domain/             → Доменный слой: сущности, value objects, интерфейсы репозиториев
  usecase/            → Слой бизнес-логики: use cases (application services)
  repo/postgres/      → Инфраструктурный слой: реализация репозиториев (PostgreSQL + sqlx)
  handler/            → Транспортный слой: HTTP handlers, middleware, DTO, роутер (Chi)
  config/             → Конфигурация приложения
pkg/
  jwt/                → Генерация и валидация JWT токенов (access + refresh)
  hash/               → Хэширование паролей (bcrypt)
  httputil/           → HTTP-утилиты (стандартизированные ответы)
migrations/           → SQL миграции (goose v3)
tests/integration/    → Интеграционные тесты (testcontainers)
```

**Зависимости направлены внутрь:** handler → usecase → domain ← repo/postgres.
Доменный слой не зависит ни от чего внешнего.

## Стек технологий

| Компонент       | Технология                     |
|-----------------|--------------------------------|
| Язык            | Go 1.25                        |
| HTTP Router     | go-chi/chi v5                  |
| База данных     | PostgreSQL 16                  |
| SQL             | jmoiron/sqlx (raw SQL)         |
| Миграции        | pressly/goose v3               |
| Аутентификация  | golang-jwt/jwt v5 (access + refresh) |
| Хэширование     | golang.org/x/crypto/bcrypt     |
| Валидация       | go-playground/validator v10    |
| Конфигурация    | caarlos0/env v11               |
| Логирование     | log/slog (stdlib)              |
| Swagger         | swaggo/swag                    |
| Тесты           | testify + testcontainers-go    |
| Контейнеризация | Docker + docker-compose        |

## Модели данных

### Пользователь (User)

| Поле            | Тип          | Описание                  |
|-----------------|--------------|---------------------------|
| id              | UUID         | Уникальный идентификатор  |
| name            | string       | Имя пользователя          |
| email           | string       | Email (уникальный)        |
| password_hash   | string       | Хэш пароля (bcrypt)       |
| created_at      | timestamp    | Дата регистрации          |

### Задача (Task)

| Поле            | Тип          | Описание                  |
|-----------------|--------------|---------------------------|
| id              | UUID         | Уникальный идентификатор  |
| title           | string       | Заголовок (1-255 символов)|
| description     | string       | Описание (опционально)    |
| status          | enum         | new / in_progress / done  |
| deadline        | timestamp    | Срок выполнения (опц.)    |
| created_at      | timestamp    | Дата создания             |
| updated_at      | timestamp    | Дата обновления           |
| deleted_at      | timestamp    | Soft delete (NULL = active)|
| author_id       | UUID (FK)    | ID автора (создателя)     |

## API Endpoints

### Аутентификация

| Метод  | Путь                | Описание               | Auth |
|--------|---------------------|------------------------|------|
| POST   | /api/v1/auth/register | Регистрация           | -    |
| POST   | /api/v1/auth/login    | Вход (получение JWT)  | -    |
| POST   | /api/v1/auth/refresh  | Обновление токенов    | Refresh |

### Задачи

| Метод  | Путь                    | Описание                    | Auth    |
|--------|-------------------------|-----------------------------|---------|
| POST   | /api/v1/tasks           | Создать задачу              | Bearer  |
| GET    | /api/v1/tasks           | Список задач (с фильтрами) | Bearer  |
| GET    | /api/v1/tasks/:id       | Получить задачу по ID       | Bearer  |
| PUT    | /api/v1/tasks/:id       | Обновить задачу             | Bearer  |
| DELETE | /api/v1/tasks/:id       | Удалить задачу (soft)       | Bearer  |
| GET    | /api/v1/tasks/stats     | Статистика по задачам       | Bearer  |

### Фильтрация и пагинация (GET /api/v1/tasks)

| Параметр         | Тип     | Пример                          | Описание                         |
|------------------|---------|----------------------------------|----------------------------------|
| status           | string  | ?status=done                     | Фильтр по статусу               |
| search           | string  | ?search=deploy                   | Поиск по заголовку (ILIKE)       |
| overdue          | bool    | ?overdue=true                    | Только просроченные              |
| deadline_before  | date    | ?deadline_before=2026-04-01      | Дедлайн до указанной даты        |
| deadline_after   | date    | ?deadline_after=2026-03-01       | Дедлайн после указанной даты     |
| created_after    | date    | ?created_after=2026-01-01        | Созданные после даты             |
| created_before   | date    | ?created_before=2026-12-31       | Созданные до даты                |
| sort_by          | string  | ?sort_by=deadline                | Поле сортировки (created_at, deadline, status, title) |
| order            | string  | ?order=asc                       | Направление (asc/desc)           |
| page             | int     | ?page=1                          | Номер страницы (от 1)            |
| limit            | int     | ?limit=20                        | Элементов на страницу (max 100)  |

### Формат ответов

Успешный ответ:
```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 42
  }
}
```

Ответ с ошибкой:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      {"field": "email", "message": "must be a valid email"}
    ]
  }
}
```

### Примеры запросов

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "email": "john@example.com", "password": "secret123"}'
```

**Создание задачи:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Deploy v2", "description": "Deploy to prod", "deadline": "2026-04-01T12:00:00Z"}'
```

**Список с фильтрами:**
```bash
curl "http://localhost:8080/api/v1/tasks?status=in_progress&overdue=true&sort_by=deadline&order=asc&page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

**Статистика:**
```bash
curl http://localhost:8080/api/v1/tasks/stats \
  -H "Authorization: Bearer <token>"
```

Ответ:
```json
{
  "data": {
    "total": 25,
    "by_status": {
      "new": 10,
      "in_progress": 8,
      "done": 7
    },
    "overdue": 3
  }
}
```

## Запуск

### С Docker Compose (рекомендуется)

```bash
cp .env.example .env
docker-compose up --build
```

API доступно на `http://localhost:8080`.
Swagger UI: `http://localhost:8080/swagger/`.

### Локально

Требования: Go 1.25+, PostgreSQL 16+.

```bash
cp .env.example .env
# отредактировать .env с параметрами вашей БД

just migrate-up   # применить миграции
just run          # запустить сервер
```

### Полезные команды

```bash
just build          # собрать бинарник
just test           # юнит-тесты
just test-integ     # интеграционные тесты (нужен Docker)
just lint           # линтер (golangci-lint)
just swagger        # сгенерировать Swagger docs
just migrate-up     # применить миграции
just migrate-down   # откатить последнюю миграцию
just migrate-create # создать новую миграцию
```

## Разработка

Проект разрабатывается по методологии **TDD (Test-Driven Development)**:
1. Написать тест, описывающий желаемое поведение
2. Убедиться, что тест падает (red)
3. Написать минимальный код для прохождения теста (green)
4. Отрефакторить (refactor)

Подход **RDD (Readme-Driven Development)**:
этот README написан до первой строки кода и служит спецификацией проекта.
Весь API, модели данных и архитектура описаны здесь до начала реализации.

## Лицензия

MIT
