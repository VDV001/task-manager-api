package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/daniilgit/task-manager-api/pkg/httputil"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(tokens usecase.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
				return
			}

			bearer, token, found := strings.Cut(header, " ")
			if !found || strings.ToLower(bearer) != "bearer" || token == "" {
				httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
				return
			}

			userID, err := tokens.ParseAccessUserID(token)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var ErrNoUserID = errors.New("user ID not found in context")

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, ErrNoUserID
	}
	return id, nil
}
