package dto

import "library/internal/domain"

type CreateAuthorInput struct {
	domain.Author
}

type CreateAuthorOutput struct {
	Id int `json:"id"`
}
