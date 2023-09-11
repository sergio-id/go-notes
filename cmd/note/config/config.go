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
		App        configs.App   `yaml:"app"`
		HTTP       configs.HTTP  `yaml:"http"`
		Logger     logger.Config `yaml:"logger"`
		PG         PG            `yaml:"postgres"`
		AuthClient AuthClient    `yaml:"auth_client"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DsnURL  string `env-required:"true" yaml:"dsn_url" env:"PG_DSN_URL"`
	}

	AuthClient struct {
		URL string `env-required:"true" yaml:"url" env:"AUTH_CLIENT_URL"`
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
