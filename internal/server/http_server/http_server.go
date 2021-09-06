package http_server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain"
)

// ServerHTTP
type ServerHTTP struct {
	router     *chi.Mux
	httpServer *http.Server

	domain *domain.Domain
}

func New(url string, d *domain.Domain) (*ServerHTTP, error) {
	server := ServerHTTP{}
	serverWrapper := ServerInterfaceWrapper{
		Handler: &server,
	}

	// Set middlewares and handlers
	r := chi.NewRouter()
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(hlog.AccessHandler(mwAccessLogger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		// Auth
		r.Group(func(r chi.Router) {
			r.Use(mwAuthTokenValidatorAndSetter)

			// User
			r.Group(func(r chi.Router) {
				r.Use(mwUserIDSetter)

				r.Get("/users", serverWrapper.GetUsers)
				r.Get("/users/{UserID}", serverWrapper.GetUsersUserID)
				r.Put("/users/{UserID}", serverWrapper.PutUsersUserID)
				r.Delete("/users/{UserID}", serverWrapper.DeleteUsersUserID)
			})
		})

		// Noauth
		r.Group(func(r chi.Router) {
			// Token
			r.Post("/tokens", serverWrapper.PostTokens)

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
