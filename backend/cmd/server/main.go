package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/data"
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/mailer"
	"portfolio-backend/internal/middleware"
)

func main() {
	cfg := config.Load()

	store, err := data.Load()
	if err != nil {
		log.Fatalf("failed to load content data: %v", err)
	}

	var m mailer.Mailer
	if cfg.SMTPHost != "" && cfg.ContactTo != "" {
		m = mailer.NewSMTP(cfg)
	} else {
		log.Println("warning: SMTP not configured (SMTP_HOST/CONTACT_TO), contact form will only log messages")
		m = mailer.ConsoleMailer{}
	}

	h := handlers.New(store, m)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.Health)
	mux.HandleFunc("GET /api/profile", h.Profile)
	mux.HandleFunc("GET /api/skills", h.Skills)
	mux.HandleFunc("GET /api/projects", h.Projects)
	mux.HandleFunc("GET /api/experience", h.Experience)

	contactLimiter := middleware.RateLimit(cfg.RateLimitPerMin, time.Minute)
	mux.Handle("POST /api/contact", contactLimiter(http.HandlerFunc(h.Contact)))

	handler := middleware.Chain(mux, middleware.Logging, middleware.CORS(cfg.AllowedOrigins))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
}
