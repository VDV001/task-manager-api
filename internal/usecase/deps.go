package usecase

import "github.com/google/uuid"

// PasswordHasher — абстракция для хэширования паролей (инвертированная зависимость).
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// TokenPair — пара JWT токенов (access + refresh).
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenManager — абстракция для управления JWT токенами.
type TokenManager interface {
	GeneratePair(userID uuid.UUID) (*TokenPair, error)
	ParseAccessUserID(token string) (uuid.UUID, error)
	ParseRefreshUserID(token string) (uuid.UUID, error)
}
