package dto

import "library/internal/domain"

type CreateBookInput struct {
	domain.Book
}

type CreateBookOutput struct {
	Id int `json:"id"`
}
