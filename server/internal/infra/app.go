package infra

import (
	"context"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/server/internal/services/music_server"
)

type Service interface {
	Run() error
	Stop() error
}

type App struct {
	Services []Service
}

func New() (*App, error) {
	ca, err := config.New()
	if err != nil {
		return nil, err
	}

	ms, err := music_server.New(ca)
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
		s.Run()
	}
}

func (a *App) Stop() {
	for _, s := range a.Services {
		s.Stop()
	}
}
