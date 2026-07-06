package main

import (
	"context"
	"errors"
	"library/internal/api"
	"library/internal/db"
	"library/internal/repository"
	"library/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := db.InitDb("library.db")
	if err != nil {
		log.Fatalf("initializing db issues: %v", err)
	}
	defer db.Close()

	store := repository.NewLibraryStore(db)
	svc := service.NewLibraryService(store)
	handler := api.NewHandler(svc)

	router := handler.SetupRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

		<-signalChan
		log.Println("graceful shutdown started")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown: %v", err)
		}
		log.Println("server is down")
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listening issues: %v", err)
	}
}
