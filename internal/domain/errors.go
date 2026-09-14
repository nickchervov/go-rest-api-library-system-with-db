package domain

type LibraryError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e LibraryError) Error() string {
	return e.Message
}

var (
	ErrIncorrectId               = &LibraryError{Code: 400, Message: "incorrect id"}
	ErrBookNotFound              = &LibraryError{Code: 404, Message: "book not found"}
	ErrAuthorNotFound            = &LibraryError{Code: 404, Message: "author not found"}
	ErrReaderNotFound            = &LibraryError{Code: 404, Message: "reader not found"}
	ErrBorrowingNotFound         = &LibraryError{Code: 404, Message: "borrowing not found"}
	ErrAuthorHasRelatedBook      = &LibraryError{Code: 400, Message: "cant delete author: has books"}
	ErrCreateBookAuthorNotExists = &LibraryError{Code: 400, Message: "cant create book: not found author"}
	ErrBookHasRelatedBorrowing   = &LibraryError{Code: 400, Message: "cant delete book: has borrowing"}
	ErrReaderHasRelatedBorrowing = &LibraryError{Code: 400, Message: "cant delete reader: has borrowing"}
	ErrNotFoundBookOrReader      = &LibraryError{Code: 404, Message: "cant add borrowing: not found book or reader"}
	ErrBookAlreadyOccupied       = &LibraryError{Code: 400, Message: "book is already occupied"}
	ErrValidation                = &LibraryError{Code: 400, Message: "incorrect data"}
)
