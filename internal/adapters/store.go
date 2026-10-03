package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"library/internal/domain"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// Только всё связанное с SQL

type Sqlite struct {
	db *sql.DB
}

func migrateUp(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("create db instance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://pkg/migrations", "sqlite", driver)
	if err != nil {
		return fmt.Errorf("create migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up migration: %w", err)
	}
	return nil
}

func NewLibraryStore(dbPath string) (*Sqlite, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open connection issues: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection issues: %w", err)
	}

	if err := migrateUp(db); err != nil {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	return &Sqlite{db: db}, nil
}

// Methods for authors table
func (s *Sqlite) Close() {
	s.db.Close()
}

func (s *Sqlite) GetAllAuthors() ([]domain.Author, error) {
	rows, err := s.db.Query("SELECT id, name, country FROM authors")
	if err != nil {
		return nil, fmt.Errorf("get all authors issues: %w", err)
	}
	defer rows.Close()

	authors := make([]domain.Author, 0)
	for rows.Next() {
		var a domain.Author
		if err := rows.Scan(&a.Id, &a.Name, &a.Country); err != nil {
			return nil, fmt.Errorf("scanning all authors rows issues: %w", err)
		}
		authors = append(authors, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("after scanning all authors rows issues: %w", err)
	}

	return authors, nil
}

func (s *Sqlite) GetAuthorById(id int) (domain.Author, error) {
	row := s.db.QueryRow("SELECT id, name, country FROM authors WHERE id = :id", sql.Named("id", id))

	var author domain.Author
	if err := row.Scan(&author.Id, &author.Name, &author.Country); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Author{}, domain.ErrAuthorNotFound
		}
		return domain.Author{}, fmt.Errorf("scanning author by id row issues: %w", err)
	}
	return author, nil
}

func (s *Sqlite) CreateAuthor(a domain.Author) (int, error) {
	res, err := s.db.Exec("INSERT INTO authors (name, country) VALUES (:name,:country)",
		sql.Named("name", a.Name),
		sql.Named("country", a.Country))
	if err != nil {
		return 0, fmt.Errorf("create author issues: %w", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get author last inserted id issues: %w", err)
	}

	return int(lastId), nil
}

func (s *Sqlite) UpdateAuthor(id int, a domain.Author) error {
	res, err := s.db.Exec("UPDATE authors SET name = :name, country = :country WHERE id = :id",
		sql.Named("name", a.Name),
		sql.Named("country", a.Country),
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("update author issues: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return domain.ErrAuthorNotFound
	}
	return nil
}

func (s *Sqlite) DeleteAuthor(id int) error {
	var bookExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE authors_id = :id)", sql.Named("id", id)).Scan(&bookExists)
	if err != nil {
		return fmt.Errorf("check author's books issues: %w", err)
	}
	if bookExists {
		return domain.ErrAuthorHasRelatedBook
	}
	res, err := s.db.Exec("DELETE FROM authors WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete author issues: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected issues: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrAuthorNotFound
	}
	return nil
}

// Methods for books table
func (s *Sqlite) GetAllBooks() ([]domain.BookResponse, error) {
	rows, err := s.db.Query(`SELECT b.id, b.title, b.isbn, b.year,a.id, a.name, a.country FROM books b
		JOIN authors a ON b.authors_id = a.id`)
	if err != nil {
		return nil, fmt.Errorf("get all books issues: %w", err)
	}
	defer rows.Close()

	books := make([]domain.BookResponse, 0)
	for rows.Next() {
		var b domain.BookResponse
		if err := rows.Scan(&b.Id, &b.Title, &b.Isbn, &b.Year, &b.Author.Id, &b.Author.Name, &b.Author.Country); err != nil {
			return nil, fmt.Errorf("scanning all books rows issues: %w", err)
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("after scanning all books rows issues: %w", err)
	}

	return books, nil
}

func (s *Sqlite) GetBookById(id int) (domain.BookResponse, error) {
	row := s.db.QueryRow(`SELECT b.id, b.title, b.isbn, b.year,a.id, a.name, a.country FROM books b
		JOIN authors a ON b.authors_id = a.id WHERE b.id = :id`, sql.Named("id", id))

	var book domain.BookResponse
	if err := row.Scan(&book.Id, &book.Title, &book.Isbn, &book.Year, &book.Author.Id, &book.Author.Name, &book.Author.Country); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.BookResponse{}, domain.ErrBookNotFound
		}
		return domain.BookResponse{}, fmt.Errorf("scanning book by id row issues: %w", err)
	}
	return book, nil
}

func (s *Sqlite) CreateBook(b domain.Book) (int, error) {
	var authorExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM authors WHERE id = :id)", sql.Named("id", b.AuthorsId)).Scan(&authorExists)
	if err != nil {
		return 0, fmt.Errorf("check existing author issues: %w", err)
	}
	if !authorExists {
		return 0, domain.ErrCreateBookAuthorNotExists
	}
	res, err := s.db.Exec("INSERT INTO books (title, isbn, year, authors_id) VALUES (:title, :isbn, :year, :authors_id)",
		sql.Named("title", b.Title),
		sql.Named("isbn", b.Isbn),
		sql.Named("year", b.Year),
		sql.Named("authors_id", b.AuthorsId))
	if err != nil {
		return 0, fmt.Errorf("create book issues: %w", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get book last inserted id issues: %w", err)
	}

	return int(lastId), nil
}

func (s *Sqlite) UpdateBook(id int, b domain.Book) error {
	res, err := s.db.Exec("UPDATE books SET title = :title, isbn = :isbn, year = :year, authors_id = :authors_id WHERE id = :id",
		sql.Named("title", b.Title),
		sql.Named("isbn", b.Isbn),
		sql.Named("year", b.Year),
		sql.Named("authors_id", b.AuthorsId),
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("update book issues: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

func (s *Sqlite) DeleteBook(id int) error {
	var borrowingsExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM borrowing WHERE book_id = :id)", sql.Named("id", id)).Scan(&borrowingsExists)
	if err != nil {
		return fmt.Errorf("check book's borrowings issues: %w", err)
	}
	if borrowingsExists {
		return domain.ErrBookHasRelatedBorrowing
	}
	res, err := s.db.Exec("DELETE FROM books WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete book issues: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected issues: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

// Methods for readers table
func (s *Sqlite) GetAllReaders() ([]domain.Reader, error) {
	rows, err := s.db.Query("SELECT id, name, email, phone FROM readers")
	if err != nil {
		return nil, fmt.Errorf("get all readers issues: %w", err)
	}
	defer rows.Close()

	readers := make([]domain.Reader, 0)
	for rows.Next() {
		var r domain.Reader
		if err := rows.Scan(&r.Id, &r.Name, &r.Email, &r.Phone); err != nil {
			return nil, fmt.Errorf("scanning readers rows issues: %w", err)
		}
		readers = append(readers, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("after scanning all reader rows issues: %w", err)
	}
	return readers, nil
}

func (s *Sqlite) CreateReader(r domain.Reader) (int, error) {
	res, err := s.db.Exec("INSERT INTO readers (name,email,phone) VALUES (:name,:email,:phone)",
		sql.Named("name", r.Name),
		sql.Named("email", r.Email),
		sql.Named("phone", r.Phone))
	if err != nil {
		return 0, fmt.Errorf("create reader issues: %w", err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get reader last inserted id issues: %w", err)
	}
	return int(lastId), nil
}

func (s *Sqlite) DeleteReader(id int) error {
	var borrowingsExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM borrowing WHERE readers_id = :id)", sql.Named("id", id)).Scan(&borrowingsExists)
	if err != nil {
		return fmt.Errorf("check reader's borrowings issues: %w", err)
	}
	if borrowingsExists {
		return domain.ErrReaderHasRelatedBorrowing
	}
	res, err := s.db.Exec("DELETE FROM readers WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete reader issues: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected issues: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrReaderNotFound
	}
	return nil
}

// Methods for borrowing table
func (s *Sqlite) GetActiveBorrowings() ([]domain.BorrowingResponse, error) {
	rows, err := s.db.Query(`
	SELECT 
		b.id, 
		book.id, 
		book.title, 
		book.isbn, 
		book.year,
		a.id, 
		a.name, 
		a.country,
		r.id,
		r.name,
		r.email,
		r.phone,
		b.borrow_date,
		b.return_date
	FROM borrowing b
	JOIN books book ON b.book_id = book.id
	JOIN authors a ON book.authors_id = a.id
	JOIN readers r ON b.readers_id = r.id
	WHERE b.return_date IS NULL`)
	if err != nil {
		return nil, fmt.Errorf("get active borrowing issues: %w", err)
	}
	defer rows.Close()

	borrowings := make([]domain.BorrowingResponse, 0)
	for rows.Next() {
		var b domain.BorrowingResponse
		if err := rows.Scan(&b.Id, &b.Book.Id, &b.Book.Title, &b.Book.Isbn, &b.Book.Year, &b.Book.Author.Id, &b.Book.Author.Name,
			&b.Book.Author.Country, &b.Reader.Id, &b.Reader.Name, &b.Reader.Email, &b.Reader.Phone, &b.BorrowDate, &b.ReturnDate); err != nil {
			return nil, fmt.Errorf("scanning borrowings rows issues: %w", err)
		}
		borrowings = append(borrowings, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("after scanning all reader rows issues: %w", err)
	}
	return borrowings, nil
}

func (s *Sqlite) TakeOffBook(b domain.Borrowing) (int, error) {
	var bookReaderExists bool
	err := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM books WHERE id = :book_id) AND EXISTS (SELECT 1 FROM readers WHERE id = :reader_id)",
		sql.Named("book_id", b.BookId), sql.Named("reader_id", b.ReadersId)).Scan(&bookReaderExists)
	if err != nil {
		return 0, fmt.Errorf("checking book and reader existing issues: %w", err)
	}
	if !bookReaderExists {
		return 0, domain.ErrNotFoundBookOrReader
	}
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx issues: %w", err)
	}
	defer tx.Rollback()

	var dbBookID int
	var activeBorrowID sql.NullInt64
	err = tx.QueryRow(`SELECT b.id, br.id FROM books b LEFT JOIN borrowing br ON b.id = br.book_id AND br.return_date IS NULL
		WHERE b.id = :id`, sql.Named("id", b.BookId)).Scan(&dbBookID, &activeBorrowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrBookNotFound
		}
		return 0, fmt.Errorf("select book issues: %w", err)
	}

	if activeBorrowID.Valid {
		return 0, domain.ErrBookAlreadyOccupied
	}

	res, err := tx.Exec("INSERT INTO borrowing (book_id, readers_id, borrow_date) VALUES (:book_id, :readers_id, :borrow_date)",
		sql.Named("book_id", b.BookId),
		sql.Named("readers_id", b.ReadersId),
		sql.Named("borrow_date", b.BorrowDate))
	if err != nil {
		return 0, fmt.Errorf("create borrowing issues: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit issues: %w", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get borrowing last inserted id issues: %w", err)
	}
	return int(lastId), nil
}

func (s *Sqlite) ReturnBook(idBorrowing int, date string) error {
	res, err := s.db.Exec("UPDATE borrowing SET return_date = :date WHERE id = :id", sql.Named("date", date), sql.Named("id", idBorrowing))
	if err != nil {
		return fmt.Errorf("return book issues: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return domain.ErrBorrowingNotFound
	}
	return nil
}
