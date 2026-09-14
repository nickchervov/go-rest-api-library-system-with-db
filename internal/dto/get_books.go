package dto

import "library/internal/domain"

type GetBooksOutput struct {
	Books []domain.BookResponse `json:"books"`
}
