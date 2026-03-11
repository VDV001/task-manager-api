package handler

import (
	"errors"
	"net/http"

	"github.com/daniilgit/task-manager-api/internal/domain"
	"github.com/daniilgit/task-manager-api/internal/usecase"
	"github.com/daniilgit/task-manager-api/pkg/httputil"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	auth     *usecase.AuthUseCase
	validate *validator.Validate
}

func NewAuthHandler(auth *usecase.AuthUseCase, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{auth: auth, validate: validate}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body RegisterRequest true "Данные регистрации"
// @Success      201 {object} httputil.Response{data=TokenResponse}
// @Failure      409 {object} httputil.ErrorResponse
// @Failure      422 {object} httputil.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if errs := validateStruct(h.validate, req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	pair, err := h.auth.Register(r.Context(), usecase.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "CONFLICT", "email already registered")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to register user")
		return
	}

	httputil.Success(w, http.StatusCreated, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

// Login godoc
// @Summary      Вход в систему
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body LoginRequest true "Данные входа"
// @Success      200 {object} httputil.Response{data=TokenResponse}
// @Failure      401 {object} httputil.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if errs := validateStruct(h.validate, req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	pair, err := h.auth.Login(r.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid email or password")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to login")
		return
	}

	httputil.Success(w, http.StatusOK, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

// Refresh godoc
// @Summary      Обновление токенов
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body RefreshRequest true "Refresh token"
// @Success      200 {object} httputil.Response{data=TokenResponse}
// @Failure      401 {object} httputil.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := decodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if errs := validateStruct(h.validate, req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	pair, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			httputil.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired refresh token")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to refresh tokens")
		return
	}

	httputil.Success(w, http.StatusOK, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func validateStruct(v *validator.Validate, s any) []httputil.FieldError {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return []httputil.FieldError{{Field: "body", Message: "validation failed"}}
	}

	result := make([]httputil.FieldError, len(validationErrs))
	for i, fe := range validationErrs {
		result[i] = httputil.FieldError{
			Field:   fe.Field(),
			Message: validationMessage(fe),
		}
	}
	return result
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "oneof":
		return "must be one of: " + fe.Param()
	default:
		return "invalid value"
	}
}
