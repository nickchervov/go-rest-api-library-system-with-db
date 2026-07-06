package service

import (
	"fmt"
	"library/internal/models"
	"library/internal/repository"
	"library/internal/validator"
)

// Только всё связанное с бизнес логикой (работы сервера)

type LibraryService struct {
	store repository.Storage
}

func NewLibraryService(store *repository.Store) *LibraryService {
	return &LibraryService{store: store}
}

func (s *LibraryService) GetAllAuthors() ([]models.Author, error) {
	return s.store.GetAllAuthors()
}

func (s *LibraryService) GetAuthorById(id int) (models.Author, error) {
	if id <= 0 {
		return models.Author{}, fmt.Errorf("id must be more then zero")
	}
	return s.store.GetAuthorById(id)
}

func (s *LibraryService) CreateAuthor(a models.Author) (int, error) {
	if err := validator.Validator.Struct(a); err != nil {
		return 0, fmt.Errorf("validation author issues: %w", err)
	}
	return s.store.CreateAuthor(a)
}

func (s *LibraryService) UpdateAuthor(id int, a models.Author) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	if err := validator.Validator.Struct(a); err != nil {
		return fmt.Errorf("validation author issues: %w", err)
	}
	return s.store.UpdateAuthor(id, a)
}

func (s *LibraryService) DeleteAuthor(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	return s.store.DeleteAuthor(id)
}

func (s *LibraryService) GetAllBooks() ([]models.BookResponse, error) {
	return s.store.GetAllBooks()
}

func (s *LibraryService) GetBookById(id int) (models.BookResponse, error) {
	if id <= 0 {
		return models.BookResponse{}, fmt.Errorf("id must be more then zero")
	}
	return s.store.GetBookById(id)
}

func (s *LibraryService) CreateBook(b models.Book) (int, error) {
	if err := validator.Validator.Struct(b); err != nil {
		return 0, fmt.Errorf("validation book issues: %w", err)
	}
	return s.store.CreateBook(b)
}

func (s *LibraryService) UpdateBook(id int, b models.Book) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	if err := validator.Validator.Struct(b); err != nil {
		return fmt.Errorf("validation book issues: %w", err)
	}
	return s.store.UpdateBook(id, b)
}

func (s *LibraryService) DeleteBook(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	return s.store.DeleteBook(id)
}

func (s *LibraryService) GetAllReaders() ([]models.Reader, error) {
	return s.store.GetAllReaders()
}

func (s *LibraryService) CreateReader(r models.Reader) (int, error) {
	if err := validator.Validator.Struct(r); err != nil {
		return 0, fmt.Errorf("validation reader issues: %w", err)
	}
	return s.store.CreateReader(r)
}

func (s *LibraryService) DeleteReader(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	return s.store.DeleteReader(id)
}

func (s *LibraryService) GetActiveBorrowings() ([]models.BorrowingResponse, error) {
	return s.store.GetActiveBorrowings()
}

func (s *LibraryService) TakeOffBook(br models.Borrowing) (int, error) {
	if err := validator.Validator.Struct(br); err != nil {
		return 0, fmt.Errorf("validation borrowing issues: %w", err)
	}
	return s.store.TakeOffBook(br)
}

func (s *LibraryService) ReturnBook(id int, date models.Date) error {
	if id <= 0 {
		return fmt.Errorf("id must be more then zero")
	}
	if err := validator.Validator.Struct(date); err != nil {
		return fmt.Errorf("validation date issues: %w", err)
	}
	return s.store.ReturnBook(id, date.ReturnDate)
}
