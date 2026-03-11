package httputil

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details []FieldError   `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, Response{Data: data})
}

func SuccessWithMeta(w http.ResponseWriter, status int, data any, meta any) {
	JSON(w, status, Response{Data: data, Meta: meta})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationError(w http.ResponseWriter, details []FieldError) {
	JSON(w, http.StatusUnprocessableEntity, ErrorResponse{
		Error: ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input",
			Details: details,
		},
	})
}
