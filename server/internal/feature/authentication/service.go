package authentication

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/class-nonetype/kin/server/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorEmailTaken         = errors.New("account with this email already exists.")
	ErrorUsernameTaken      = errors.New("account with this username already exists.")
	ErrorInvalidCredentials = errors.New("invalid username or password.")
)

type Service struct {
	q          *store.Queries
	sessionTTL time.Duration
}

func NewService(q *store.Queries, sessionTTL time.Duration) *Service {
	return &Service{q: q, sessionTTL: sessionTTL}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (store.User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return store.User{}, err
	}

	user, err := s.q.CreateUser(ctx, store.CreateUserParams{
		Username: username,
		Email:    email,
		Password: hash,
	})

	var pgErr *pgconn.PgError

	// 23505 es el codigo de postgresql para violacion de unicidad
	// constraintname dice cual se violo:
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return store.User{}, ErrorEmailTaken
		case "users_username_key":
			return store.User{}, ErrorUsernameTaken
		}
	}
	if err != nil {
		return store.User{}, fmt.Errorf("%w", err)
	}

	return user, nil
}

var dummyHash, _ = HashPassword("dummy-password")

func (s *Service) SignIn(ctx context.Context, email, password string) (string, store.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		VerifyPassword(password, dummyHash)
		return "", store.User{}, ErrorInvalidCredentials
	}
	if err != nil {
		return "", store.User{}, fmt.Errorf("%w", err)
	}

	ok, err := VerifyPassword(password, user.Password)
	if err != nil {
		return "", store.User{}, err
	}
	if !ok {
		return "", store.User{}, ErrorInvalidCredentials
	}

	token, hash := newToken()
	err = s.q.CreateSession(ctx, store.CreateSessionParams{
		Token:     hash,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.sessionTTL),
	})
	if err != nil {
		return "", store.User{}, fmt.Errorf("%w", err)
	}

	return token, user, nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (store.GetSessionUserRow, error) {
	user, err := s.q.GetSessionUser(ctx, hashToken(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return store.GetSessionUserRow{}, ErrorInvalidCredentials
	}
	if err != nil {
		return store.GetSessionUserRow{}, fmt.Errorf("%w", err)
	}
	return user, nil
}

func (s *Service) SignOut(ctx context.Context, token string) error {
	return s.q.DeleteSession(ctx, hashToken(token))
}
