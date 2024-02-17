package storage

import (
	"context"

	"github.com/izzdalfk/app-skeleton/go/internal/shared/logger"
)

type dummyStorage struct{}

func NewDummy() *dummyStorage {
	return &dummyStorage{}
}

func (s *dummyStorage) GetEntity(ctx context.Context, entity string) string {
	// test that logger can retrieve from any service components
	svcLogger := logger.GetFromContext(ctx)
	if svcLogger != nil {
		svcLogger.Info("GetEntity called from dummyStorage", logger.Field{
			Key:   "entity",
			Value: entity,
		})
	}

	return entity
}
