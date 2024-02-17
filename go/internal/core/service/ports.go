package service

import "context"

type Storage interface {
	GetEntity(ctx context.Context, entity string) string
}
