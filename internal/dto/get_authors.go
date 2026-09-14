package dto

import "library/internal/domain"

type GetAuthorsOutput struct {
	Authors []domain.Author `json:"authors"`
}
