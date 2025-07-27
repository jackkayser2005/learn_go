package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jackkayser2005/learn_go/internal/handlers"
	"github.com/jackkayser2005/learn_go/internal/store"
)

func main() {
	st := store.NewMemStore()
	h := handlers.Handlers{Store: st}

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	r.Post("/shorten", h.Shorten)
	r.Get("/{code}", h.Redirect)

	srv := &http.Server{
		Addr:              ":8000",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// start
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()
	log.Println("listening on :8000")

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("bye")

}
