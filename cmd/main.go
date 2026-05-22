package main

import (
	"library/internal/api"
	"library/internal/db"
	"library/internal/repository"
	"library/internal/service"
	"log"
	"net/http"
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

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("opening connection issues: %v", err)
	}
}
