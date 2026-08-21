package db

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("Entity not found")
	ErrDuplicate     = errors.New("Entity already exists")
	ErrInvalidEntity = errors.New("Invalid entity")
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetById(ctx context.Context, id int) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*T, error)
	Count(ctx context.Context) (int64, error)
}
