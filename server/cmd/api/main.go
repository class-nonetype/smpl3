package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/class-nonetype/kin/server/internal/config"
	"github.com/class-nonetype/kin/server/internal/database"
	"github.com/class-nonetype/kin/server/internal/feature/authentication"
	"github.com/class-nonetype/kin/server/internal/feature/post"
	"github.com/class-nonetype/kin/server/internal/httpx"
	"github.com/class-nonetype/kin/server/internal/middleware"
	"github.com/class-nonetype/kin/server/internal/store"
)

func main() {

	// ctrl + c manda SIGTERM, el ctx se cancela
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newLogger(env string) *slog.Logger {
	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func run(ctx context.Context) error {
	// variables de entorno
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// logger
	logger := newLogger(cfg.Env)
	slog.SetDefault(logger)

	// base de datos
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	// cierra las conexiones cuando ru ntermina, incluso si sale por un error
	defer pool.Close()

	queries := store.New(pool)

	authService := authentication.NewService(queries, 30*24*time.Hour)
	authHandler := authentication.NewHandler(authService, cfg.Env == "production")
	postHandler := post.NewHandler(queries)

	// servidor

	router := httpx.NewRouter()
	router.Handle("GET", "/healthz", func(w http.ResponseWriter, r *http.Request) error {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return nil
	})

	api := router.Group("/api/v1")

	auth := api.Group("/authentication")
	auth.Handle("POST", "/sign-up", authHandler.SignUp)
	auth.Handle("POST", "/sign-in", authHandler.SignIn)
	auth.Handle("POST", "/sign-out", authHandler.SignOut)
	auth.Group("", authHandler.RequireUser).Handle("GET", "/me", authHandler.Me)

	posts := api.Group("/posts")
	posts.Handle("GET", "", postHandler.List)
	posts.Handle("GET", "/{id}", postHandler.Get)

	privatePosts := posts.Group("", authHandler.RequireUser)
	privatePosts.Handle("POST", "", postHandler.Create)
	privatePosts.Handle("DELETE", "/{id}", postHandler.Delete)

	handler := middleware.Recover(middleware.Logger(middleware.CORS(cfg.AllowedOrigins)(router)))

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", cfg.Addr, err)
	}

	for _, route := range router.Routes() {
		logger.Info("route", "endpoint", route)
	}

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	if ip := net.ParseIP(host); host == "" || (ip != nil && ip.IsUnspecified()) {
		host = "localhost"
	}
	logger.Info("server listening", "host", host, "port", port, "url", "http://"+net.JoinHostPort(host, port))

	errCh := make(chan error, 1)

	// el servidor corre en una goroutine
	// asi run queda libre para esperar en el select lo que pase primero:
	// que el servidor falle o que llega la señal
	go func() {
		if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// deja de aceptar conexiones nuevas y espera hasta 10 segundos a que terminen
	// las peticiones en curso
	return srv.Shutdown(shutdownCtx)
}
