package dto

import "library/internal/domain"

type TakeOffBookInput struct {
	domain.Borrowing
}

type TakeOffBookOutput struct {
	Id int `json:"id"`
}
