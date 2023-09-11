package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	configs "github.com/sergio-id/go-notes/pkg/config"
	"github.com/sergio-id/go-notes/pkg/logger"
	"os"
	"path/filepath"
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

	pathToConfigFile, err := getPathToConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getPathToConfig")
	}

	err = cleanenv.ReadConfig(pathToConfigFile, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cleanenv.ReadConfig")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cleanenv.ReadEnv")
	}

	return cfg, nil
}

func getPathToConfig() (string, error) {
	dir, err := getPath()
	if err != nil {
		return "", errors.Wrap(err, "getPath")
	}

	return filepath.Join(dir, "config.yml"), nil
}

/**
 * getPath returns the path to the output directory (for local development)
 *
 * input: ../go-grpc-notes-microservices-test
 * output: ../go-grpc-notes-microservices-test/cmd/user.
 */
func getPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "os.Getwd")
	}

	if dir == "/" {
		return dir, nil
	}

	outputDir := filepath.Join(dir, "cmd", "user")

	return outputDir, nil
}
