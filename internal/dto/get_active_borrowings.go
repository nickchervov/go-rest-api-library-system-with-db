package dto

import "library/internal/domain"

type GetActiveBorrowingsOutput struct {
	Borrowings []domain.BorrowingResponse
}
