package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/class-nonetype/kin/server/internal/httpx"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				slog.ErrorContext(r.Context(), "panic",
					"panic", v,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				httpx.WriteJSON(w, http.StatusInternalServerError,
					map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
