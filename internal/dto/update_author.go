package dto

import "library/internal/domain"

type UpdateAuthorInput struct {
	Id        int
	NewAuthor domain.Author
}
