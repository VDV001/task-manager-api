package jwt

import (
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/google/uuid"
)

var _ usecase.TokenManager = (*ManagerAdapter)(nil)

// ManagerAdapter адаптирует Manager к интерфейсу usecase.TokenManager.
type ManagerAdapter struct {
	m *Manager
}

func NewManagerAdapter(m *Manager) *ManagerAdapter {
	return &ManagerAdapter{m: m}
}

func (a *ManagerAdapter) GeneratePair(userID uuid.UUID) (*usecase.TokenPair, error) {
	access, refresh, err := a.m.GeneratePair(userID)
	if err != nil {
		return nil, err
	}
	return &usecase.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (a *ManagerAdapter) ParseAccessUserID(token string) (uuid.UUID, error) {
	claims, err := a.m.ParseAccess(token)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func (a *ManagerAdapter) ParseRefreshUserID(token string) (uuid.UUID, error) {
	claims, err := a.m.ParseRefresh(token)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}
