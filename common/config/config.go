package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	_path = "./config.yaml"
)

type Config struct {
	MusicServer MusicServerConfig `yaml:"music_server"`
	Listener    MusicServerConfig `yaml:"listener"`
	Redis       RedisConfig       `yaml:"redis"`
}

type MusicServerConfig struct {
	Address string `yaml:"address"`
}

type ListenerConfig struct {
	Address string `yaml:"address"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func New() (*Config, error) {
	var ca Config

	if err := cleanenv.ReadConfig(_path, &ca); err != nil {
		return nil, err
	}

	return &ca, nil
}
