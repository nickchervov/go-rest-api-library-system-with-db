package models

type Author struct {
	Id      int    `json:"id" validate:"omitempty"`
	Name    string `json:"name" validate:"required,alphaspace"`
	Country string `json:"country" validate:"required,alphaspace"`
}

type Reader struct {
	Id    int    `json:"id" validate:"omitempty"`
	Name  string `json:"name" validate:"required,alphaspace"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"omitempty"`
}

type Book struct {
	Id        int    `json:"id" validate:"omitempty"`
	Title     string `json:"title" validate:"required"`
	Isbn      string `json:"isbn" validate:"required,numericspecs"`
	Year      int    `json:"year" validate:"required,numeric"`
	AuthorsId int    `json:"authors_id" validate:"required,numeric"`
}

type BookResponse struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Isbn   string `json:"isbn"`
	Year   int    `json:"year"`
	Author Author `json:"authors"`
}

type Borrowing struct {
	Id         int    `json:"id" validate:"omitempty"`
	BookId     int    `json:"book_id" validate:"required,numeric"`
	ReadersId  int    `json:"readers_id" validate:"required,numeric"`
	BorrowDate string `json:"borrow_date" validate:"required,datetime=2006-01-02"`
	ReturnDate string `json:"return_date" validate:"omitempty"`
}
type BorrowingResponse struct {
	Id         int          `json:"id"`
	Book       BookResponse `json:"book"`
	Reader     Reader       `json:"reader"`
	BorrowDate string       `json:"borrow_date"`
	ReturnDate *string      `json:"return_date"`
}

type Date struct {
	ReturnDate string `json:"return_date" validate:"required,datetime=2006-01-02"`
}
