package dto

import "library/internal/domain"

type UpdateBookInput struct {
	Id      int
	NewBook domain.Book
}
