package post

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/class-nonetype/kin/server/internal/feature/authentication"
	"github.com/class-nonetype/kin/server/internal/httpx"
	"github.com/class-nonetype/kin/server/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	defaultPageSize = 20
	maxPageSize     = 50
)

type Handler struct {
	q *store.Queries
}

func NewHandler(q *store.Queries) *Handler {
	return &Handler{q: q}
}

type author struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type postResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    author    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type listResponse struct {
	Items      []postResponse `json:"items"`
	NextCursor *uuid.UUID     `json:"nextCursor,omitempty"`
}

func toResponse(p store.GetPostRow) postResponse {
	return postResponse{
		ID:        p.ID,
		Title:     p.Title,
		Body:      p.Body,
		Author:    author{ID: p.AuthorID, Username: p.AuthorUsername},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type createRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (req *createRequest) validate() error {
	req.Title = strings.TrimSpace(req.Title)

	switch n := utf8.RuneCountInString(req.Title); {
	case n == 0:
		return httpx.NewError(http.StatusUnprocessableEntity, "title is required")
	case n > 200:
		return httpx.NewError(http.StatusUnprocessableEntity, "title must be at most 200 characters")
	case utf8.RuneCountInString(req.Body) > 20000:
		return httpx.NewError(http.StatusUnprocessableEntity, "body must be at most 20000 characters")
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.UserFrom(r.Context())
	if !ok {
		return httpx.NewError(http.StatusUnauthorized, "authentication required")
	}

	var req createRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		return httpx.NewError(http.StatusBadRequest, "invalid JSON")
	}
	if err := req.validate(); err != nil {
		return err
	}

	id, err := h.q.CreatePost(r.Context(), store.CreatePostParams{
		AuthorID: user.ID,
		Title:    req.Title,
		Body:     req.Body,
	})
	if err != nil {
		return err
	}

	post, err := h.q.GetPost(r.Context(), id)
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusCreated, toResponse(post))
	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httpx.NewError(http.StatusNotFound, "post not found")
	}

	post, err := h.q.GetPost(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		return httpx.NewError(http.StatusNotFound, "post not found")
	}
	if err != nil {
		return err
	}

	httpx.WriteJSON(w, http.StatusOK, toResponse(post))
	return nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	pageSize := defaultPageSize
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return httpx.NewError(http.StatusBadRequest, "limit must be a positive integer")
		}
		pageSize = min(n, maxPageSize)
	}

	var before uuid.NullUUID
	if raw := r.URL.Query().Get("before"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return httpx.NewError(http.StatusBadRequest, "invalid cursor")
		}
		before = uuid.NullUUID{UUID: id, Valid: true}
	}

	rows, err := h.q.ListPosts(r.Context(), store.ListPostsParams{
		Before:   before,
		PageSize: int32(pageSize + 1),
	})
	if err != nil {
		return err
	}

	resp := listResponse{Items: make([]postResponse, 0, len(rows))}
	if len(rows) > pageSize {
		rows = rows[:pageSize]
		resp.NextCursor = &rows[pageSize-1].ID
	}
	for _, row := range rows {
		resp.Items = append(resp.Items, toResponse(store.GetPostRow(row)))
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.UserFrom(r.Context())
	if !ok {
		return httpx.NewError(http.StatusUnauthorized, "authentication required")
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httpx.NewError(http.StatusNotFound, "post not found")
	}

	deleted, err := h.q.DeletePost(r.Context(), store.DeletePostParams{ID: id, AuthorID: user.ID})
	if err != nil {
		return err
	}
	if deleted == 0 {
		return httpx.NewError(http.StatusNotFound, "post not found")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
