package infra

import (
	"context"

	"github.com/kuroko-shirai/together/client/internal/services/listener"
	"github.com/kuroko-shirai/together/common/config"
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

	ms, err := listener.New(ca)
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
