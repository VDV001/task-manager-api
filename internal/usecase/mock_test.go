package usecase_test

import (
	"context"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/google/uuid"
)

// --- UserRepository mock ---

type mockUserRepo struct {
	createFn     func(ctx context.Context, user *domain.User) error
	getByIDFn    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	if m.createFn == nil {
		panic("mockUserRepo.Create: unexpected call")
	}
	return m.createFn(ctx, user)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getByIDFn == nil {
		panic("mockUserRepo.GetByID: unexpected call")
	}
	return m.getByIDFn(ctx, id)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFn == nil {
		panic("mockUserRepo.GetByEmail: unexpected call")
	}
	return m.getByEmailFn(ctx, email)
}

// --- TaskRepository mock ---

type mockTaskRepo struct {
	createFn     func(ctx context.Context, task *domain.Task) error
	getByIDFn    func(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	updateFn     func(ctx context.Context, task *domain.Task) error
	softDeleteFn func(ctx context.Context, id uuid.UUID) error
	listFn       func(ctx context.Context, authorID uuid.UUID, filter *domain.TaskFilter) ([]domain.Task, int, error)
	getStatsFn   func(ctx context.Context, authorID uuid.UUID) (*domain.TaskStats, error)
}

func (m *mockTaskRepo) Create(ctx context.Context, task *domain.Task) error {
	if m.createFn == nil {
		panic("mockTaskRepo.Create: unexpected call")
	}
	return m.createFn(ctx, task)
}

func (m *mockTaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	if m.getByIDFn == nil {
		panic("mockTaskRepo.GetByID: unexpected call")
	}
	return m.getByIDFn(ctx, id)
}

func (m *mockTaskRepo) Update(ctx context.Context, task *domain.Task) error {
	if m.updateFn == nil {
		panic("mockTaskRepo.Update: unexpected call")
	}
	return m.updateFn(ctx, task)
}

func (m *mockTaskRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if m.softDeleteFn == nil {
		panic("mockTaskRepo.SoftDelete: unexpected call")
	}
	return m.softDeleteFn(ctx, id)
}

func (m *mockTaskRepo) List(ctx context.Context, authorID uuid.UUID, filter *domain.TaskFilter) ([]domain.Task, int, error) {
	if m.listFn == nil {
		panic("mockTaskRepo.List: unexpected call")
	}
	return m.listFn(ctx, authorID, filter)
}

func (m *mockTaskRepo) GetStats(ctx context.Context, authorID uuid.UUID) (*domain.TaskStats, error) {
	if m.getStatsFn == nil {
		panic("mockTaskRepo.GetStats: unexpected call")
	}
	return m.getStatsFn(ctx, authorID)
}

// --- PasswordHasher mock ---

type mockHasher struct {
	hashFn    func(password string) (string, error)
	compareFn func(hash, password string) bool
}

func (m *mockHasher) Hash(password string) (string, error) {
	if m.hashFn == nil {
		panic("mockHasher.Hash: unexpected call")
	}
	return m.hashFn(password)
}

func (m *mockHasher) Compare(hash, password string) bool {
	if m.compareFn == nil {
		panic("mockHasher.Compare: unexpected call")
	}
	return m.compareFn(hash, password)
}

// --- TokenManager mock ---

type mockTokenManager struct {
	generatePairFn       func(userID uuid.UUID) (*usecase.TokenPair, error)
	parseAccessUserIDFn  func(token string) (uuid.UUID, error)
	parseRefreshUserIDFn func(token string) (uuid.UUID, error)
}

func (m *mockTokenManager) GeneratePair(userID uuid.UUID) (*usecase.TokenPair, error) {
	if m.generatePairFn == nil {
		panic("mockTokenManager.GeneratePair: unexpected call")
	}
	return m.generatePairFn(userID)
}

func (m *mockTokenManager) ParseAccessUserID(token string) (uuid.UUID, error) {
	if m.parseAccessUserIDFn == nil {
		panic("mockTokenManager.ParseAccessUserID: unexpected call")
	}
	return m.parseAccessUserIDFn(token)
}

func (m *mockTokenManager) ParseRefreshUserID(token string) (uuid.UUID, error) {
	if m.parseRefreshUserIDFn == nil {
		panic("mockTokenManager.ParseRefreshUserID: unexpected call")
	}
	return m.parseRefreshUserIDFn(token)
}
