package service

import (
	"fmt"
	"library/internal/domain"
	"library/internal/dto"
	"library/pkg/validator"
	"time"
)

// Только всё связанное с бизнес логикой (работы сервера)

type Storage interface {
	Close()
	GetAllAuthors() ([]domain.Author, error)
	GetAuthorById(id int) (domain.Author, error)
	CreateAuthor(a domain.Author) (int, error)
	UpdateAuthor(id int, a domain.Author) error
	DeleteAuthor(id int) error
	GetAllBooks() ([]domain.BookResponse, error)
	GetBookById(id int) (domain.BookResponse, error)
	CreateBook(b domain.Book) (int, error)
	UpdateBook(id int, b domain.Book) error
	DeleteBook(id int) error
	GetAllReaders() ([]domain.Reader, error)
	CreateReader(r domain.Reader) (int, error)
	DeleteReader(id int) error
	GetActiveBorrowings() ([]domain.BorrowingResponse, error)
	TakeOffBook(b domain.Borrowing) (int, error)
	ReturnBook(idBorrowing int, date string) error
}

type LibraryService struct {
	store Storage
}

func NewLibraryService(store Storage) *LibraryService {
	return &LibraryService{store: store}
}

func (s *LibraryService) GetAllAuthors() (dto.GetAuthorsOutput, error) {
	authors, err := s.store.GetAllAuthors()
	if err != nil {
		return dto.GetAuthorsOutput{}, fmt.Errorf("get all authors: %w", err)
	}

	return dto.GetAuthorsOutput{Authors: authors}, nil
}

func (s *LibraryService) GetAuthorById(input dto.GetAuthorInput) (dto.GetAuthorOutput, error) {
	if input.Id <= 0 {
		return dto.GetAuthorOutput{}, domain.ErrIncorrectId
	}
	author, err := s.store.GetAuthorById(input.Id)
	if err != nil {
		return dto.GetAuthorOutput{}, fmt.Errorf("get author by %d: %w", input.Id, err)
	}

	return dto.GetAuthorOutput{Author: author}, nil
}

func (s *LibraryService) CreateAuthor(input dto.CreateAuthorInput) (dto.CreateAuthorOutput, error) {
	if err := validator.Validator.Struct(input); err != nil {
		return dto.CreateAuthorOutput{}, &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	id, err := s.store.CreateAuthor(input.Author)
	if err != nil {
		return dto.CreateAuthorOutput{}, fmt.Errorf("create author: %w", err)
	}
	return dto.CreateAuthorOutput{Id: id}, nil
}

func (s *LibraryService) UpdateAuthor(input dto.UpdateAuthorInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := validator.Validator.Struct(input.NewAuthor); err != nil {
		return &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	if err := s.store.UpdateAuthor(input.Id, input.NewAuthor); err != nil {
		return fmt.Errorf("update author: %w", err)
	}
	return nil
}

func (s *LibraryService) DeleteAuthor(input dto.DeleteAuthorInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := s.store.DeleteAuthor(input.Id); err != nil {
		return fmt.Errorf("delete author: %w", err)
	}
	return nil
}

func (s *LibraryService) GetAllBooks() (dto.GetBooksOutput, error) {
	books, err := s.store.GetAllBooks()
	if err != nil {
		return dto.GetBooksOutput{}, fmt.Errorf("get all books: %w", err)
	}
	return dto.GetBooksOutput{Books: books}, nil
}

func (s *LibraryService) GetBookById(input dto.GetBookInput) (dto.GetBookOutput, error) {
	if input.Id <= 0 {
		return dto.GetBookOutput{}, domain.ErrIncorrectId
	}
	book, err := s.store.GetBookById(input.Id)
	if err != nil {
		return dto.GetBookOutput{}, fmt.Errorf("get book by %d: %w", input.Id, err)
	}
	return dto.GetBookOutput{BookResponse: book}, nil
}

func (s *LibraryService) CreateBook(input dto.CreateBookInput) (dto.CreateBookOutput, error) {
	if err := validator.Validator.Struct(input); err != nil {
		return dto.CreateBookOutput{}, &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	id, err := s.store.CreateBook(input.Book)
	if err != nil {
		return dto.CreateBookOutput{}, fmt.Errorf("create book: %w", err)
	}
	return dto.CreateBookOutput{Id: id}, nil
}

func (s *LibraryService) UpdateBook(input dto.UpdateBookInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := validator.Validator.Struct(input.NewBook); err != nil {
		return &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	if err := s.store.UpdateBook(input.Id, input.NewBook); err != nil {
		return fmt.Errorf("update book: %w", err)
	}
	return nil
}

func (s *LibraryService) DeleteBook(input dto.DeleteBookInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := s.store.DeleteBook(input.Id); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	return nil
}

func (s *LibraryService) GetAllReaders() (dto.GetReadersOutput, error) {
	readers, err := s.store.GetAllReaders()
	if err != nil {
		return dto.GetReadersOutput{}, fmt.Errorf("get readers: %w", err)
	}
	return dto.GetReadersOutput{Readers: readers}, nil
}

func (s *LibraryService) CreateReader(input dto.CreateReaderInput) (dto.CreateReaderOutput, error) {
	if err := validator.Validator.Struct(input); err != nil {
		return dto.CreateReaderOutput{}, &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	id, err := s.store.CreateReader(input.Reader)
	if err != nil {
		return dto.CreateReaderOutput{}, fmt.Errorf("create reader: %w", err)
	}
	return dto.CreateReaderOutput{Id: id}, nil
}

func (s *LibraryService) DeleteReader(input dto.DeleteReaderInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := s.store.DeleteReader(input.Id); err != nil {
		return fmt.Errorf("delete reader: %w", err)
	}
	return nil
}

func (s *LibraryService) GetActiveBorrowings() (dto.GetActiveBorrowingsOutput, error) {
	borrowings, err := s.store.GetActiveBorrowings()
	if err != nil {
		return dto.GetActiveBorrowingsOutput{}, fmt.Errorf("get active borrowing: %w", err)
	}
	return dto.GetActiveBorrowingsOutput{Borrowings: borrowings}, nil
}

func (s *LibraryService) TakeOffBook(input dto.TakeOffBookInput) (dto.TakeOffBookOutput, error) {
	if err := validator.Validator.Struct(input); err != nil {
		return dto.TakeOffBookOutput{}, &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	input.BorrowDate = time.Now().Format("2006-01-02")
	id, err := s.store.TakeOffBook(input.Borrowing)
	if err != nil {
		return dto.TakeOffBookOutput{}, fmt.Errorf("take off book: %w", err)
	}
	return dto.TakeOffBookOutput{Id: id}, nil
}

func (s *LibraryService) ReturnBook(input dto.ReturnBookInput) error {
	if input.Id <= 0 {
		return domain.ErrIncorrectId
	}
	if err := validator.Validator.Struct(input.Date); err != nil {
		return &domain.LibraryError{Code: 400, Message: err.Error()}
	}
	if err := s.store.ReturnBook(input.Id, input.Date.ReturnDate); err != nil {
		return fmt.Errorf("return book: %w", err)
	}
	return nil
}
