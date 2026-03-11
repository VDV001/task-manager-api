package usecase_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAuthUseCase(
	userRepo domain.UserRepository,
	hasher usecase.PasswordHasher,
	tokens usecase.TokenManager,
) *usecase.AuthUseCase {
	return usecase.NewAuthUseCase(userRepo, hasher, tokens, slog.Default())
}

func TestAuthUseCase_Register(t *testing.T) {
	t.Parallel()

	validInput := usecase.RegisterInput{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	}

	tests := []struct {
		name      string
		input     usecase.RegisterInput
		setupRepo func() *mockUserRepo
		wantErr   error
	}{
		{
			name:  "success",
			input: validInput,
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					createFn: func(_ context.Context, user *domain.User) error {
						assert.Equal(t, "John", user.Name)
						assert.Equal(t, "john@example.com", user.Email)
						assert.NotEmpty(t, user.ID)
						return nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name:  "email already exists",
			input: validInput,
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					createFn: func(_ context.Context, _ *domain.User) error {
						return domain.ErrAlreadyExists
					},
				}
			},
			wantErr: domain.ErrAlreadyExists,
		},
		{
			name:  "repo error on create",
			input: validInput,
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					createFn: func(_ context.Context, _ *domain.User) error {
						return errors.New("db down")
					},
				}
			},
			wantErr: errors.New("create user: db down"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repo := tt.setupRepo()
			hasher := &mockHasher{
				hashFn: func(_ string) (string, error) {
					return "$2a$hashed", nil
				},
			}
			tokens := &mockTokenManager{
				generatePairFn: func(_ uuid.UUID) (*usecase.TokenPair, error) {
					return &usecase.TokenPair{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					}, nil
				},
			}

			uc := newAuthUseCase(repo, hasher, tokens)
			pair, err := uc.Register(ctx, tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, domain.ErrAlreadyExists) {
					assert.ErrorIs(t, err, domain.ErrAlreadyExists)
				} else {
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
				assert.Nil(t, pair)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "access-token", pair.AccessToken)
				assert.Equal(t, "refresh-token", pair.RefreshToken)
			}
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name      string
		input     usecase.LoginInput
		setupRepo func() *mockUserRepo
		compareFn func(hash, password string) bool
		wantErr   error
	}{
		{
			name:  "success",
			input: usecase.LoginInput{Email: "john@example.com", Password: "secret123"},
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return &domain.User{ID: userID, PasswordHash: "$2a$hashed"}, nil
					},
				}
			},
			compareFn: func(_, _ string) bool { return true },
			wantErr:   nil,
		},
		{
			name:  "user not found",
			input: usecase.LoginInput{Email: "unknown@example.com", Password: "secret123"},
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return nil, domain.ErrNotFound
					},
				}
			},
			compareFn: func(_, _ string) bool { return true },
			wantErr:   domain.ErrInvalidCredentials,
		},
		{
			name:  "db error propagates",
			input: usecase.LoginInput{Email: "john@example.com", Password: "secret123"},
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return nil, errors.New("connection refused")
					},
				}
			},
			compareFn: func(_, _ string) bool { return true },
			wantErr:   errors.New("get user by email: connection refused"),
		},
		{
			name:  "wrong password",
			input: usecase.LoginInput{Email: "john@example.com", Password: "wrong"},
			setupRepo: func() *mockUserRepo {
				return &mockUserRepo{
					getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
						return &domain.User{ID: userID, PasswordHash: "$2a$hashed"}, nil
					},
				}
			},
			compareFn: func(_, _ string) bool { return false },
			wantErr:   domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repo := tt.setupRepo()
			hasher := &mockHasher{
				compareFn: tt.compareFn,
			}
			tokens := &mockTokenManager{
				generatePairFn: func(_ uuid.UUID) (*usecase.TokenPair, error) {
					return &usecase.TokenPair{AccessToken: "at", RefreshToken: "rt"}, nil
				},
			}

			uc := newAuthUseCase(repo, hasher, tokens)
			pair, err := uc.Login(ctx, tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				if errors.Is(tt.wantErr, domain.ErrInvalidCredentials) {
					assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
				} else {
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
				assert.Nil(t, pair)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, pair.AccessToken)
			}
		})
	}
}

func TestAuthUseCase_Refresh(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name           string
		refreshToken   string
		parseRefreshFn func(token string) (uuid.UUID, error)
		getByIDFn      func(ctx context.Context, id uuid.UUID) (*domain.User, error)
		wantErr        error
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh",
			parseRefreshFn: func(_ string) (uuid.UUID, error) {
				return userID, nil
			},
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
				return &domain.User{ID: userID}, nil
			},
			wantErr: nil,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "bad-token",
			parseRefreshFn: func(_ string) (uuid.UUID, error) {
				return uuid.Nil, errors.New("invalid")
			},
			getByIDFn: nil,
			wantErr:   domain.ErrTokenInvalid,
		},
		{
			name:         "user deleted after token issued",
			refreshToken: "valid-refresh",
			parseRefreshFn: func(_ string) (uuid.UUID, error) {
				return userID, nil
			},
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.User, error) {
				return nil, domain.ErrNotFound
			},
			wantErr: domain.ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repo := &mockUserRepo{
				getByIDFn: tt.getByIDFn,
			}
			tokens := &mockTokenManager{
				parseRefreshUserIDFn: tt.parseRefreshFn,
				generatePairFn: func(_ uuid.UUID) (*usecase.TokenPair, error) {
					return &usecase.TokenPair{AccessToken: "new-at", RefreshToken: "new-rt"}, nil
				},
			}

			uc := newAuthUseCase(repo, &mockHasher{}, tokens)
			pair, err := uc.Refresh(ctx, tt.refreshToken)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "new-at", pair.AccessToken)
			}
		})
	}
}
