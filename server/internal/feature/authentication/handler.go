package authentication

import (
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/class-nonetype/kin/server/internal/httpx"
	"github.com/google/uuid"
)

const sessionCookie = "kin_session"

type Handler struct {
	svc          *Service
	secureCookie bool
}

func NewHandler(svc *Service, secureCookie bool) *Handler {
	return &Handler{svc: svc, secureCookie: secureCookie}
}

var usernameRe = regexp.MustCompile(`^[a-z0-9_]{3,30}$`)

type signUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *signUpRequest) validate() error {
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Email = strings.TrimSpace(req.Email)

	switch {
	case !usernameRe.MatchString(req.Username):
		return httpx.NewError(http.StatusUnprocessableEntity, "username must be 3 to 30 characters and contain only letters, numbers, and underscores")

	case len(req.Email) > 254:
		return httpx.NewError(http.StatusUnprocessableEntity, "email is too long")

	case !isEmail(req.Email):
		return httpx.NewError(http.StatusUnprocessableEntity, "invalid email")

	case len(req.Password) < 8 || len(req.Password) > 72:
		return httpx.NewError(http.StatusUnprocessableEntity, "password must be between 8 and 72 characters")

	}
	return nil
}

func isEmail(s string) bool {
	addr, err := mail.ParseAddress(s)
	return err == nil && addr.Address == s
}

type userResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var req signUpRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		return httpx.NewError(http.StatusBadRequest, "JSON inválido")
	}
	if err := req.validate(); err != nil {
		return err
	}

	user, err := h.svc.Register(r.Context(), req.Username, req.Email, req.Password)
	switch {
	case errors.Is(err, ErrorEmailTaken), errors.Is(err, ErrorUsernameTaken):
		return httpx.NewError(http.StatusConflict, err.Error())
	case err != nil:
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	return nil
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var req signInRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		return httpx.NewError(http.StatusBadRequest, "invalid JSON")
	}

	token, user, err := h.svc.SignIn(r.Context(), strings.TrimSpace(req.Email), req.Password)
	if errors.Is(err, ErrorInvalidCredentials) {
		return httpx.NewError(http.StatusUnauthorized, err.Error())
	}
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.svc.sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	httpx.WriteJSON(w, http.StatusOK, userResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	return nil
}
