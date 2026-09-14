package connectors

import (
	"encoding/json"
	"errors"
	"library/internal/domain"
	"library/internal/dto"
	"library/internal/service"
	"library/pkg/render"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *service.LibraryService
}

func NewHandler(svc *service.LibraryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAuthors(w http.ResponseWriter, r *http.Request) {
	output, err := h.svc.GetAllAuthors()
	if err != nil {
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get all authors internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) GetAuthorById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id: " + err.Error()})
		return
	}
	input := dto.GetAuthorInput{Id: id}
	output, err := h.svc.GetAuthorById(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get author internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAuthorInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body " + err.Error()})
		return
	}
	output, err := h.svc.CreateAuthor(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "create author internal error"})
		return
	}
	render.JSONResponse(w, http.StatusCreated, output)
}

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	var author domain.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{})
		return
	}
	input := dto.UpdateAuthorInput{
		Id:        id,
		NewAuthor: author,
	}
	if err := h.svc.UpdateAuthor(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "update author internal error"})
		return
	}
	author.Id = id
	render.JSONResponse(w, http.StatusOK, author)
}

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	input := dto.DeleteAuthorInput{
		Id: id,
	}
	if err := h.svc.DeleteAuthor(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "delete author internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	output, err := h.svc.GetAllBooks()
	if err != nil {
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get books internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) GetBookById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	input := dto.GetBookInput{Id: id}
	output, err := h.svc.GetBookById(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get book by id internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateBookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body " + err.Error()})
		return
	}

	output, err := h.svc.CreateBook(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "create book internal error"})
		return
	}
	render.JSONResponse(w, http.StatusCreated, output)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body " + err.Error()})
		return
	}
	input := dto.UpdateBookInput{
		Id:      id,
		NewBook: book,
	}
	if err := h.svc.UpdateBook(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "update book internal error"})
		return
	}
	book.Id = id
	render.JSONResponse(w, http.StatusOK, book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	input := dto.DeleteBookInput{Id: id}
	if err := h.svc.DeleteBook(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "delete book internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetReaders(w http.ResponseWriter, r *http.Request) {
	output, err := h.svc.GetAllReaders()
	if err != nil {
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get readers internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) CreateReader(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateReaderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body " + err.Error()})
		return
	}
	output, err := h.svc.CreateReader(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "create reader internal error"})
		return
	}
	render.JSONResponse(w, http.StatusCreated, output)
}

func (h *Handler) DeleteReader(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	input := dto.DeleteReaderInput{Id: id}
	if err := h.svc.DeleteReader(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "delete reader internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetActiveBorrowings(w http.ResponseWriter, r *http.Request) {
	output, err := h.svc.GetActiveBorrowings()
	if err != nil {
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "get active borrowing internal error"})
		return
	}
	render.JSONResponse(w, http.StatusOK, output)
}

func (h *Handler) TakeOffBook(w http.ResponseWriter, r *http.Request) {
	var input dto.TakeOffBookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body: " + err.Error()})
		return
	}
	output, err := h.svc.TakeOffBook(input)
	if err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "take of book internal error"})
		return
	}
	render.JSONResponse(w, http.StatusCreated, output)
}

func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "incorrect id"})
		return
	}
	var date domain.Date
	if err := json.NewDecoder(r.Body).Decode(&date); err != nil {
		render.JSONResponse(w, http.StatusBadRequest, map[string]string{"message": "decoding request body: " + err.Error()})
		return
	}
	input := dto.ReturnBookInput{
		Id:   id,
		Date: date,
	}
	if err := h.svc.ReturnBook(input); err != nil {
		var targetErr *domain.LibraryError
		if errors.As(err, &targetErr) {
			render.JSONResponse(w, targetErr.Code, targetErr)
			return
		}
		log.Println(err)
		render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "return book internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
