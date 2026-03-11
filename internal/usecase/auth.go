package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthUseCase struct {
	users  domain.UserRepository
	hasher PasswordHasher
	tokens TokenManager
	log    *slog.Logger
}

func NewAuthUseCase(
	users domain.UserRepository,
	hasher PasswordHasher,
	tokens TokenManager,
	log *slog.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		users:  users,
		hasher: hasher,
		tokens: tokens,
		log:    log,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*TokenPair, error) {
	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	if err := uc.users.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	uc.log.Info("user registered", slog.String("user_id", user.ID.String()))

	pair, err := uc.tokens.GeneratePair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return pair, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*TokenPair, error) {
	user, err := uc.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if !uc.hasher.Compare(user.PasswordHash, input.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	pair, err := uc.tokens.GeneratePair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return pair, nil
}

func (uc *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userID, err := uc.tokens.ParseRefreshUserID(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	// Проверяем, что пользователь всё ещё существует.
	if _, err := uc.users.GetByID(ctx, userID); err != nil {
		return nil, domain.ErrTokenInvalid
	}

	pair, err := uc.tokens.GeneratePair(userID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return pair, nil
}
