package dto

import "library/internal/domain"

type CreateReaderInput struct {
	domain.Reader
}

type CreateReaderOutput struct {
	Id int `json:"id"`
}
