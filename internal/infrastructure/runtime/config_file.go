package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type fileConfig struct {
	Database struct {
		Driver string `yaml:"driver"`
		Path   string `yaml:"path"`
	} `yaml:"database"`
	Auth struct {
		PIN                   string `yaml:"pin"`
		AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
		RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
		PermanentTokens       *bool  `yaml:"permanent_tokens"`
	} `yaml:"auth"`
}

func loadConfigFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read config file: %w", err)
	}

	var cfg fileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	if cfg.Database.Driver != "" {
		_ = os.Setenv("DB_DRIVER", cfg.Database.Driver)
	}
	if cfg.Database.Path != "" {
		_ = os.Setenv("DB_SCHEMA", cfg.Database.Path)
	}
	if cfg.Auth.PIN != "" {
		_ = os.Setenv("AUTH_PIN", cfg.Auth.PIN)
	}
	if cfg.Auth.AccessTokenSecretKey != "" {
		_ = os.Setenv("AUTH_ACCESS_TOKEN_SECRET_KEY", cfg.Auth.AccessTokenSecretKey)
	}
	if cfg.Auth.RefreshTokenSecretKey != "" {
		_ = os.Setenv("AUTH_REFRESH_TOKEN_SECRET_KEY", cfg.Auth.RefreshTokenSecretKey)
	}
	if cfg.Auth.PermanentTokens != nil {
		_ = os.Setenv("AUTH_PERMANENT_TOKENS", fmt.Sprintf("%t", *cfg.Auth.PermanentTokens))
	}

	return nil
}
