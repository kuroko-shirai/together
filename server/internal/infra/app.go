package infra

import (
	"context"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/server/internal/services/music_server"
)

type Service interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type App struct {
	Services []Service
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	ms, err := music_server.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Services: []Service{
			ms,
		},
	}, nil
}

func (a *App) Run(cxt context.Context) {
	for _, s := range a.Services {
		s.Run(cxt)
	}
}

func (a *App) Stop(cxt context.Context) {
	for _, s := range a.Services {
		s.Stop(cxt)
	}
}
