package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	_path = "./config.yaml"
)

type Config struct {
	MusicServer MusicServerConfig `yaml:"music_server"`
	Listener    ListenerConfig    `yaml:"listener"`
}

type MusicServerConfig struct {
	Address string `yaml:"address"`
}

type ListenerConfig struct {
	Address string `yaml:"address"`
}

func New() (*Config, error) {
	var ca Config

	if err := cleanenv.ReadConfig(_path, &ca); err != nil {
		return nil, err
	}

	return &ca, nil
}
