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
Доменный слой не зависит ни от чего внешнего. Все внешние зависимости инвертированы через интерфейсы.

## Стек технологий

| Компонент       | Технология                     |
|-----------------|--------------------------------|
| Язык            | Go 1.25                        |
| HTTP Router     | go-chi/chi v5                  |
| База данных     | PostgreSQL 17                  |
| SQL             | jmoiron/sqlx (raw SQL)         |
| Миграции        | pressly/goose v3               |
| Аутентификация  | golang-jwt/jwt v5 (access + refresh) |
| Хэширование     | golang.org/x/crypto/bcrypt     |
| Валидация       | go-playground/validator v10    |
| Конфигурация    | caarlos0/env v11               |
| Логирование     | log/slog (stdlib, JSON)        |
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
| description     | string       | Описание (до 10000 символов) |
| status          | enum         | new / in_progress / done  |
| deadline        | timestamp    | Срок выполнения (опц.)    |
| created_at      | timestamp    | Дата создания             |
| updated_at      | timestamp    | Дата обновления           |
| deleted_at      | timestamp    | Soft delete (NULL = active)|
| author_id       | UUID (FK)    | ID автора (создателя)     |

Для soft-deleted записей используются partial indexes в PostgreSQL — индексы строятся только по `deleted_at IS NULL`, что исключает удалённые записи из поисковых операций и экономит место.

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
| PATCH  | /api/v1/tasks/:id       | Обновить задачу (частично)  | Bearer  |
| DELETE | /api/v1/tasks/:id       | Удалить задачу (soft)       | Bearer  |
| GET    | /api/v1/tasks/stats     | Статистика по задачам       | Bearer  |

### Служебные

| Метод  | Путь      | Описание                                 | Auth |
|--------|-----------|------------------------------------------|------|
| GET    | /healthz  | Liveness probe                           | -    |
| GET    | /readyz   | Readiness probe (проверка БД)            | -    |
| GET    | /version  | Информация о сборке (версия, коммит, Go) | -    |

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

Даты принимаются в двух форматах: `YYYY-MM-DD` и `RFC3339` (`2026-04-01T12:00:00Z`).

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

Коды ошибок: `BAD_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `VALIDATION_ERROR`, `INTERNAL_ERROR`.

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

**Обновление задачи (частичное):**
```bash
curl -X PATCH http://localhost:8080/api/v1/tasks/<id> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}'
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

## Фронтенд (web/)

SPA-приложение на Vue 3 с тёмной темой и glassmorphism-дизайном.

### Стек фронтенда

| Компонент         | Технология                          |
|-------------------|-------------------------------------|
| Фреймворк         | Vue 3 (Composition API, `<script setup>`) |
| Сборка            | Vite 7                              |
| Язык              | TypeScript 5 (strict)               |
| Стили             | Tailwind CSS 4                      |
| Стейт             | Pinia 3                             |
| Роутинг           | Vue Router 4 (navigation guards)    |
| Иконки            | Lucide Vue Next                     |
| Формы             | vee-validate + Zod                  |
| HTTP-клиент       | ofetch (с interceptors)             |
| Интернационализация | vue-i18n 10 (EN / RU)            |
| Графики           | Chart.js + vue-chartjs              |
| Виртуальный скролл | @tanstack/vue-virtual              |
| Уведомления       | vue-sonner                          |
| Тестирование      | Vitest + @vue/test-utils            |
| Обратный прокси   | Caddy                               |

### Возможности

- JWT-аутентификация (access + refresh) с auto-refresh
- Защищённые маршруты (navigation guards)
- Дашборд: карточки и таблица (переключение)
- DataTable с виртуальным скроллом (10 000+ строк)
- Графики: линейный (создание задач по времени) + круговая (по статусам)
- Фильтры таблицы синхронизированы с графиками
- Пагинация (карточный вид)
- CRUD задач с модальными окнами
- Поиск, сортировка, фильтрация по статусу
- Подсветка просроченных задач
- Адаптивный дизайн (мобильные устройства)
- Демо-режим с 10K моковых задач
- Toast-уведомления
- Скелетоны загрузки
- Переключение языка (EN/RU)
- Error Boundary с retry (key-based re-mount)
- Доступность (a11y): ARIA-атрибуты, семантический HTML, навигация с клавиатуры

### Запуск фронтенда

```bash
cd web
npm install
npm run dev          # dev-сервер на http://localhost:3000
```

### Команды фронтенда

```bash
npm run dev          # запуск dev-сервера
npm run build        # TypeScript проверка + production сборка
npm run lint         # ESLint
npm run format       # Prettier
npm run test         # unit-тесты (Vitest)
npm run test:watch   # тесты в watch-режиме
npm run test:coverage # тесты с покрытием
```

### Структура фронтенда

```
web/
  src/
    api/              → HTTP-клиент и API-модули (auth, tasks)
    assets/           → Стили (Tailwind CSS)
    components/
      auth/           → AuthLayout
      charts/         → TasksLineChart, TasksPieChart
      layout/         → AppSidebar (top navbar)
      tasks/          → TaskCard, TaskFilters, TaskCreateModal, TaskEditModal, TaskPagination
      ui/             → AppButton, AppCard, AppBadge, AppInput, AppModal, AppSpinner, DataTable, AppSkeleton, ErrorBoundary
    i18n/             → Интернационализация (en.ts, ru.ts)
    layouts/          → AppLayout (основной layout с navbar)
    lib/              → Утилиты (cn, mock-data)
    pages/            → LoginPage, RegisterPage, DashboardPage, StatsPage
    router/           → Vue Router с auth guards
    stores/           → Pinia stores (auth, tasks)
    types/            → TypeScript-типы (Task, Auth, API)
    __tests__/        → Unit-тесты (Vitest)
```

## Безопасность

- **JWT**: Access + Refresh токены, HMAC-SHA256 подпись, валидация метода подписи
- **Пароли**: bcrypt хэширование (cost 12)
- **Авторизация**: Пользователь видит только свои задачи
- **Soft Delete**: Удалённые записи помечаются `deleted_at`, а не стираются физически
- **SQL Injection**: Whitelist для сортировки, параметризованные запросы
- **Request Body**: Ограничение 1 MB (MaxBytesReader)
- **Content-Type**: Проверка `application/json` на routes с body
- **CORS**: Настраиваемые origins через переменную окружения
- **Rate Limiting**: Per-IP ограничение запросов с X-RateLimit-* заголовками
- **Docker**: Non-root пользователь, healthchecks, .dockerignore

## Конфигурация

Приложение конфигурируется через переменные окружения (файл `.env`):

| Переменная        | Описание                        | По умолчанию   |
|-------------------|---------------------------------|----------------|
| SERVER_PORT       | Порт HTTP сервера               | 8080           |
| SERVER_READ_TIMEOUT  | Таймаут чтения запроса       | 10s            |
| SERVER_WRITE_TIMEOUT | Таймаут записи ответа        | 30s            |
| SERVER_SHUTDOWN_TIMEOUT | Таймаут graceful shutdown  | 15s            |
| DB_HOST           | Хост PostgreSQL                 | localhost      |
| DB_PORT           | Порт PostgreSQL                 | 5432           |
| DB_USER           | Пользователь БД                 | taskmanager    |
| DB_PASSWORD       | Пароль БД                       | taskmanager    |
| DB_NAME           | Имя БД                          | taskmanager    |
| DB_SSLMODE        | SSL режим                       | disable        |
| JWT_ACCESS_SECRET | Секрет для access токенов       | (обязательно)  |
| JWT_REFRESH_SECRET| Секрет для refresh токенов      | (обязательно)  |
| JWT_ACCESS_TTL    | Время жизни access токена       | 15m            |
| JWT_REFRESH_TTL   | Время жизни refresh токена      | 720h           |
| CORS_ORIGINS      | Разрешённые origin (через запятую) | *           |
| LOG_LEVEL         | Уровень логирования             | info           |

## Запуск

### С Docker Compose (рекомендуется)

```bash
cp .env.example .env
docker-compose up --build
```

API доступно на `http://localhost:8080`.
Swagger UI: `http://localhost:8080/swagger/`.
Фронтенд: `http://localhost` (через Caddy).

### Локально

Требования: Go 1.25+, PostgreSQL 17+.

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
just test-all       # юнит + интеграционные тесты
just lint           # линтер (golangci-lint)
just swagger        # сгенерировать Swagger docs
just migrate-up     # применить миграции
just migrate-down   # откатить последнюю миграцию
just migrate-create # создать новую миграцию
just fmt            # форматирование кода
```

## Тестирование

### Юнит-тесты

Покрывают бизнес-логику (usecase layer) с ручными моками. Table-driven тесты, запускаются параллельно.

```bash
just test
```

### Интеграционные тесты (бэкенд)

Полный flow через HTTP API с реальной PostgreSQL (testcontainers). Покрывают:
регистрация, дубликат (409), логин, CRUD задач, фильтрация, статистика, soft delete (404), авторизация (401), изоляция пользователей.

```bash
just test-integ
```

### Unit-тесты (фронтенд)

Vitest + @vue/test-utils + happy-dom. Покрывают: stores (auth, tasks), router guards, i18n (проверка ключей EN/RU, интерполяция, переключение), компоненты (AppBadge), утилиты (mock-data генератор).

```bash
cd web && npm test
```

### E2E-тесты (фронтенд)

Playwright + Chromium. Покрывают: авторизация (login/register), навигация, переключение языка, дашборд.

```bash
cd web && npx playwright test
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

Коммиты оформлены по спецификации [Conventional Commits](https://www.conventionalcommits.org/ru/v1.0.0-beta.4/).

## Лицензия

MIT
