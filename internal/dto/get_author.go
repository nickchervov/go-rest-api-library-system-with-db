package dto

import "library/internal/domain"

type GetAuthorInput struct {
	Id int `json:"id"`
}

type GetAuthorOutput struct {
	domain.Author
}
