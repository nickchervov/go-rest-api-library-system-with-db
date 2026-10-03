package connectors

import (
	"library/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(svc *service.LibraryService) http.Handler {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
	h := NewHandler(svc)
	//Authors Handlers
	r.Get("/api/authors", h.GetAuthors)
	r.Get("/api/authors/{id}", h.GetAuthorById)
	r.Post("/api/authors", h.CreateAuthor)
	r.Put("/api/authors/{id}", h.UpdateAuthor)
	r.Delete("/api/authors/{id}", h.DeleteAuthor)
	//Books Handlers
	r.Get("/api/books", h.GetBooks)
	r.Get("/api/books/{id}", h.GetBookById)
	r.Post("/api/books", h.CreateBook)
	r.Put("/api/books/{id}", h.UpdateBook)
	r.Delete("/api/books/{id}", h.DeleteBook)
	//Readers Handlers
	r.Get("/api/readers", h.GetReaders)
	r.Post("/api/readers", h.CreateReader)
	r.Delete("/api/readers/{id}", h.DeleteReader)
	//Borrowing Handlers
	r.Get("/api/borrowings/active", h.GetActiveBorrowings) //Получить все активные выдачи (где return_date IS NULL)
	r.Post("/api/borrowings", h.TakeOffBook)
	r.Put("/api/borrowings/{id}/return", h.ReturnBook)
	return r
}
