package config

import (
	"net"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kuroko-shirai/together/utils"
)

const (
	_path = "./config.yaml"
)

type Config struct {
	MusicServer MusicServerConfig `yaml:"music_server"`
	Listeners   []ListenerConfig  `yaml:"listeners"`
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

func (this *Config) GetAvailableListener() (
	listener net.Listener,
	err error,
) {
	for _, cfgListener := range this.Listeners {
		listener, err = net.Listen(
			utils.TCP,
			cfgListener.Address,
		)
		if err != nil {
			continue
		} else {
			break
		}
	}

	return listener, err
}
