package main

import (
	"context"
	"errors"
	"library/internal/connectors"
	"library/internal/repository"
	"library/internal/service"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	repo, err := repository.NewLibraryStore("library.db")
	if err != nil {
		log.Fatalf("creating repository: %v", err)
	}
	svc := service.NewLibraryService(repo)

	h := connectors.SetupRoutes(svc)

	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	go func() {
		log.Printf("server listening %s port...", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listening issues: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Starting shutdown...")
	ctxShutdown, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("shutdown issues")
	}
	log.Println("Server closed.")

	repo.Close()
}
