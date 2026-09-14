package dto

import "library/internal/domain"

type ReturnBookInput struct {
	Id   int
	Date domain.Date
}
