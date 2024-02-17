package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/izzdalfk/app-skeleton/go/internal/shared/logger"
)

type Dummy interface {
	Hello(ctx context.Context, entity string) (string, error)
}

type dummyService struct {
	storage Storage
}

type Config struct {
	Storage Storage
}

func NewDummy(config Config) *dummyService {
	return &dummyService{
		storage: config.Storage,
	}
}

func (s *dummyService) Hello(ctx context.Context, entity string) (string, error) {
	// test that logger can retrieve from any service components
	svcLogger := logger.GetFromContext(ctx)
	if svcLogger != nil {
		svcLogger.Info("Hello called from core service", logger.Field{
			Key:   "entity",
			Value: entity,
		})
	}

	// number regex used to check whether the entity contains numbers or not
	numberRegex := regexp.MustCompile("[0-9]")
	if numberRegex.MatchString(entity) {
		// this is an example how to provide specific error from core service logic
		// and later will be returned with appropriate HTTP status
		return "", ErrWrongEntity
	}

	return fmt.Sprintf("Hello %s", s.storage.GetEntity(ctx, entity)), nil
}
