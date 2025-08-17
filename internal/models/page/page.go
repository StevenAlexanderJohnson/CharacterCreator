package page

import "dndcc/internal/models"

type PageData[T any] struct {
	IsAuthenticated bool
	User            *models.Claims
	Data            T
}

func NewPageData[T any](authenticated bool, user *models.Claims, data T) *PageData[T] {
	return &PageData[T]{
		IsAuthenticated: authenticated,
		User:            user,
		Data:            data,
	}
}
