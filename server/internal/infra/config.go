package infra

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kuroko-shirai/together/server/internal/services/music_server"
)

const (
	_path = "./config.yaml"
)

type Config struct {
	MusicServer music_server.Config `yaml:"music_server"`
}

func configNew() (*Config, error) {
	var ca Config

	if err := cleanenv.ReadConfig(_path, &ca); err != nil {
		return nil, err
	}

	return &ca, nil
}
