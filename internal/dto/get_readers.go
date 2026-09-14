package dto

import "library/internal/domain"

type GetReadersOutput struct {
	Readers []domain.Reader `json:"readers"`
}
