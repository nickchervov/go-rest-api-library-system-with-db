package dto

import "library/internal/domain"

type GetBookInput struct {
	Id int
}

type GetBookOutput struct {
	domain.BookResponse
}
