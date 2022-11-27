package http_server

import (
	"context"
	"net/http"
	"time"

	"github.com/casbin/casbin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/service-auth/internal/domain"
)

// ServerHTTP
type ServerHTTP struct {
	router     *chi.Mux
	httpServer *http.Server

	domain *domain.Domain
}

func New(d *domain.Domain, url string, e *casbin.Enforcer) (*ServerHTTP, error) {
	server := ServerHTTP{}
	serverWrapper := ServerInterfaceWrapper{
		Handler: &server,
	}

	// Set middlewares
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Use(hlog.NewHandler(log.Logger))
	r.Use(mwRequestIDSetter())
	r.Use(mwOpenTracingSetter())
	r.Use(hlog.AccessHandler(mwAccessLogger))

	// Set handlers
	r.Route("/v1", func(r chi.Router) {
		// Auth
		r.Group(func(r chi.Router) {
			// Set Auth middlewares
			r.Use(mwAccessTokenValidatorAndSetter())
			r.Use(mwAuthorizer(e))

			// User
			r.Group(func(r chi.Router) {
				r.Use(mwUserIDLoggerSetter())

				r.Get("/users", serverWrapper.GetUsers)
				r.Get("/users/{UserID}", serverWrapper.GetUsersUserID)
				r.Put("/users/{UserID}", serverWrapper.PutUsersUserID)
				r.Delete("/users/{UserID}", serverWrapper.DeleteUsersUserID)
				r.Get("/users/me", serverWrapper.GetUsersMe)
				r.Put("/users/me", serverWrapper.PutUsersMe)
				r.Delete("/users/me", serverWrapper.DeleteUsersMe)
			})
		})

		// Noauth
		r.Group(func(r chi.Router) {
			// Token
			r.Post("/tokens/login", serverWrapper.PostTokensLogin)
			r.Post("/tokens/refresh", serverWrapper.PostTokensRefresh)

			// User
			r.Post("/users", serverWrapper.PostUsers)

			// Swagger
			r.Get("/", getSwaggerUIHandler(url))
			r.Get("/swagger/ui", getSwaggerUIHandler(url))
			r.Get("/swagger/spec", getSwaggerSpecHandler(url))
		})
	})

	server.domain = d
	server.router = r
	return &server, nil
}

func (s *ServerHTTP) ListenAndServe() {
	s.httpServer = &http.Server{
		Addr:    ":80",
		Handler: s.router,
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to listen and server HTTP server")
		}
	}()
}

func (s *ServerHTTP) Shutdown() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second) // Wait 5 seconds
	defer shutdownCancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown gracefully")
	}
}
