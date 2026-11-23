package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/class-nonetype/kin/server/internal/httpx"
	"github.com/class-nonetype/kin/server/internal/store"
)

type ctxKey struct{}

func UserFrom(ctx context.Context) (store.GetSessionUserRow, bool) {

	user, ok := ctx.Value(ctxKey{}).(store.GetSessionUserRow)
	return user, ok
}

func (h *Handler) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil {
			httpx.WriteError(w, r, httpx.NewError(http.StatusUnauthorized, "authentication required"))
			return
		}

		user, err := h.svc.Authenticate(r.Context(), cookie.Value)
		if errors.Is(err, ErrorInvalidCredentials) {
			httpx.WriteError(w, r, httpx.NewError(http.StatusUnauthorized, "session expired"))
			return
		}
		if err != nil {
			httpx.WriteError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// borra en los dos lados
// la sesion en la base de datos (con esto el token muere)
// y la cookie en el navegador (MaxAge: -1 le dice que la elimine ya)
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) error {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if err := h.svc.SignOut(r.Context(), cookie.Value); err != nil {
			return err
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) error {
	user, ok := UserFrom(r.Context())
	if !ok {
		return httpx.NewError(http.StatusUnauthorized, "authentication required")
	}

	httpx.WriteJSON(w, http.StatusOK, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	return nil
}
