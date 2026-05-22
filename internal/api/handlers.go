package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"library/internal/models"
	"library/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	ErrCodeInternalServer = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeConflict       = "STATUS_CONFLICT"
	ErrCodeNotFound       = "NOT_FOUND"
	CodeStatusOk          = "STATUS_OK"
	CodeStatusCreated     = "STATUS_CREATED"
	CodeDeletedOk         = "DELETED_OK"
	CodeBookReturned      = "BOOK_RETURNED"
)

type Handler struct {
	svc *service.LibraryService
}

func NewHandler(svc *service.LibraryService) *Handler {
	return &Handler{svc: svc}
}

type APIErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type APIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSONResponse(w http.ResponseWriter, status int, resp any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.svc.GetAllAuthors()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, authors)
}

func (h *Handler) GetAuthorById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id issues: %w", err).Error()})
		return
	}
	author, err := h.svc.GetAuthorById(id)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, author)
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get author from request issues: %w", err).Error()})
		return
	}
	id, err := h.svc.CreateAuthor(author)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusCreated, APIResponse{Code: CodeStatusCreated, Message: fmt.Sprint(id)})
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id issues: %w", err).Error()})
		return
	}
	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get author from request issues: %w", err).Error()})
		return
	}
	if err := h.svc.UpdateAuthor(id, author); err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	author.Id = id
	JSONResponse(w, http.StatusOK, author)
}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id issues: %w", err).Error()})
		return
	}

	if err := h.svc.DeleteAuthor(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			JSONResponse(w, http.StatusNotFound, APIErrorResponse{Code: ErrCodeNotFound, Message: err.Error()})
			return
		}
		JSONResponse(w, http.StatusConflict, APIErrorResponse{Code: ErrCodeConflict, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, APIResponse{Code: CodeDeletedOk, Message: fmt.Sprint(id)})
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.svc.GetAllBooks()
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, books)
}

func (h *Handler) GetBookById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id issues: %w", err).Error()})
		return
	}
	book, err := h.svc.GetBookById(id)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, book)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get book from request issues: %w", err).Error()})
		return
	}

	id, err := h.svc.CreateBook(book)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusCreated, APIResponse{Code: CodeStatusCreated, Message: fmt.Sprint(id)})
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id book issues: %w", err).Error()})
		return
	}
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get book from request issues: %w", err).Error()})
		return
	}
	if err := h.svc.UpdateBook(id, book); err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	book.Id = id
	JSONResponse(w, http.StatusOK, book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id book issues: %w", err).Error()})
		return
	}
	if err := h.svc.DeleteBook(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			JSONResponse(w, http.StatusNotFound, APIErrorResponse{Code: ErrCodeNotFound, Message: err.Error()})
			return
		}
		JSONResponse(w, http.StatusConflict, APIErrorResponse{Code: ErrCodeConflict, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, APIResponse{Code: CodeDeletedOk, Message: fmt.Sprint(id)})
}

func (h *Handler) GetReaders(w http.ResponseWriter, r *http.Request) {
	readers, err := h.svc.GetAllReaders()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, readers)
}

func (h *Handler) CreateReader(w http.ResponseWriter, r *http.Request) {
	var reader models.Reader
	if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get reader from request issues: %w", err).Error()})
		return
	}
	id, err := h.svc.CreateReader(reader)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusCreated, APIResponse{Code: CodeStatusCreated, Message: fmt.Sprint(id)})
}

func (h *Handler) DeleteReader(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id reader issues: %w", err).Error()})
		return
	}
	if err := h.svc.DeleteReader(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			JSONResponse(w, http.StatusNotFound, APIErrorResponse{Code: ErrCodeNotFound, Message: err.Error()})
			return
		}
		JSONResponse(w, http.StatusConflict, APIErrorResponse{Code: ErrCodeConflict, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, APIResponse{Code: CodeDeletedOk, Message: fmt.Sprint(id)})
}

func (h *Handler) GetActiveBorrowings(w http.ResponseWriter, r *http.Request) {
	activeBorrowings, err := h.svc.GetActiveBorrowings()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, activeBorrowings)
}

func (h *Handler) TakeOffBook(w http.ResponseWriter, r *http.Request) {
	var borrowing models.Borrowing
	if err := json.NewDecoder(r.Body).Decode(&borrowing); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get borrowing from request issues: %w", err).Error()})
		return
	}
	id, err := h.svc.TakeOffBook(borrowing)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusCreated, APIResponse{Code: CodeStatusCreated, Message: fmt.Sprint(id)})
}

func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get id borrowing issues: %w", err).Error()})
		return
	}
	var date models.Date
	if err := json.NewDecoder(r.Body).Decode(&date); err != nil {
		JSONResponse(w, http.StatusBadRequest, APIErrorResponse{Code: ErrCodeBadRequest, Message: fmt.Errorf("get date from request issues: %w", err).Error()})
		return
	}
	if err := h.svc.ReturnBook(id, date); err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIErrorResponse{Code: ErrCodeInternalServer, Message: err.Error()})
		return
	}
	JSONResponse(w, http.StatusOK, APIResponse{Code: CodeBookReturned, Message: fmt.Sprint(id)})
}

func (h *Handler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
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
