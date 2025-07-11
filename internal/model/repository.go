package model

import (
	"context"

	"github.com/labstack/echo"
)

const (
	CacheDuration = 2
)

type Repository[T any] interface {
	Create(context.Context, *T) (*T, *echo.HTTPError)
	Delete(context.Context, *string) *echo.HTTPError
	Get(context.Context, *string) (*T, *echo.HTTPError)
	Update(context.Context, *T) *echo.HTTPError
	GetAll(context.Context) (*[]T, *echo.HTTPError)
}

/* func TokenizeFilters(filters *string) *[]string {
	tokens := strings.Split(*filters, ",")
	return &tokens
}
*/
