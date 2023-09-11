package config

import (
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/pkg/logger"
	"github.com/sergio-id/go-notes/pkg/redis"
	"github.com/sergio-id/go-notes/pkg/security"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	configs "github.com/sergio-id/go-notes/pkg/config"
)

type (
	Config struct {
		App             configs.App     `yaml:"app"`
		HTTP            configs.HTTP    `yaml:"http"`
		Logger          logger.Config   `yaml:"logger"`
		SessionDuration int             `yaml:"session_duration"`
		UserClient      UserClient      `yaml:"user_client"`
		Redis           redis.Config    `yaml:"redis"`
		Security        security.Config `yaml:"security"`
	}

	UserClient struct {
		URL string `env-required:"true" yaml:"url" env:"USER_CLIENT_URL"`
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
