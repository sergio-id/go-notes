package config

import (
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/pkg/logger"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	configs "github.com/sergio-id/go-notes/pkg/config"
)

type (
	Config struct {
		App    configs.App   `yaml:"app"`
		HTTP   configs.HTTP  `yaml:"http"`
		Logger logger.Config `yaml:"logger"`
		GRPC   GRPC          `yaml:"client"`
	}

	GRPC struct {
		AuthURL     string `env-required:"true" yaml:"auth_client_url" env:"AUTH_CLIENT_URL"`
		NoteURL     string `env-required:"true" yaml:"note_client_url" env:"NOTE_CLIENT_URL"`
		CategoryURL string `env-required:"true" yaml:"category_client_url" env:"CATEGORY_CLIENT_URL"`
		UserURL     string `env-required:"true" yaml:"user_client_url" env:"USER_CLIENT_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "os.Getwd")
	}

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cleanenv.ReadConfig")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cleanenv.ReadEnv")
	}

	return cfg, nil
}
