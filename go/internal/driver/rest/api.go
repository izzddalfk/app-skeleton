package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/izzdalfk/app-skeleton/go/internal/core/service"
	"github.com/izzdalfk/app-skeleton/go/internal/driver/rest/middleware"
	"github.com/izzdalfk/app-skeleton/go/internal/shared/logger"
	"gopkg.in/validator.v2"
)

type API struct {
	// other API related configurations, variables, or
	// core services instance
	logger       logger.Logger
	dummyService service.Dummy
}

type APIConfig struct {
	Logger       logger.Logger `validate:"nonnil"`
	DummyService service.Dummy `validate:"nonnil"`
}

func NewAPI(config APIConfig) (*API, error) {
	err := validator.Validate(config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &API{
		logger:       config.Logger,
		dummyService: config.DummyService,
	}, nil
}

// Handler used to return `http.Handler` that required by server object
func (a *API) Handler() http.Handler {
	// NOTE: In Go v1.22, net/http standard library have router functionality
	// This codebase may updated in the future using net/http for the router
	r := chi.NewRouter()
	r.Use(
		middleware.LoggerContext(a.logger),
		middleware.LogRequest(a.logger),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, NewSuccessResp("It works!"))
	})
	r.Get("/entity/{entity}", func(w http.ResponseWriter, r *http.Request) {
		entityName := chi.URLParam(r, "entity")

		helloMessage, err := a.dummyService.Hello(r.Context(), entityName)
		if err != nil {
			render.Render(w, r, NewErrorResp(err))
			return
		}
		render.Render(w, r, NewSuccessResp(helloMessage))
	})

	return r
}
