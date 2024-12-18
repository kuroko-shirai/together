package infra

import (
	"context"

	"github.com/kuroko-shirai/together/client/internal/services/listener"
	"github.com/kuroko-shirai/together/common/config"
)

type Service interface {
	Run(context.Context) error
	Down(context.Context) error
}

type App struct {
	Services []Service
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	ms, err := listener.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Services: []Service{
			ms,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) {
	for _, s := range a.Services {
		s.Run(ctx)
	}
}

func (a *App) Down(cxt context.Context) {
	for _, s := range a.Services {
		s.Down(cxt)
	}
}
