package httpx

import (
	"errors"
	"log/slog"
	"net/http"
)

type Error struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

func NewError(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		slog.ErrorContext(r.Context(), "internal server error", "error", err, "path", r.URL.Path)
		apiErr = NewError(http.StatusInternalServerError, "internal server error")

	}

	WriteJSON(w, apiErr.Status, map[string]string{"error": apiErr.Message})
}
