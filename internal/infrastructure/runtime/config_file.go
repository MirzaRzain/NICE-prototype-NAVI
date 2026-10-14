package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type fileConfig struct {
	App struct {
		Name                 string `yaml:"name"`
		Env                  string `yaml:"env"`
		Host                 string `yaml:"host"`
		Port                 int    `yaml:"port"`
		URL                  string `yaml:"url"`
		Debug                *bool  `yaml:"debug"`
		Prometheus           *bool  `yaml:"prometheus"`
		RequestTimeout       string `yaml:"request_timeout"`
		ShutdownTimeout      string `yaml:"shutdown_timeout"`
		ShutdownWait         string `yaml:"shutdown_wait"`
		RateLimitPerSecond   int    `yaml:"rate_limit_per_second"`
		RateLimitIPPerSecond int    `yaml:"rate_limit_ip_per_second"`
		PublicAPIKey         string `yaml:"public_api_key"`
	} `yaml:"app"`
	Database struct {
		Driver             string `yaml:"driver"`
		Path               string `yaml:"path"`
		Host               string `yaml:"host"`
		Port               string `yaml:"port"`
		Username           string `yaml:"username"`
		Password           string `yaml:"password"`
		SSLMode            string `yaml:"sslmode"`
		MaxOpenConnections int    `yaml:"max_open_connections"`
		MaxIdleConnections int    `yaml:"max_idle_connections"`
		MaxLifetime        string `yaml:"max_lifetime"`
		MaxConnIdleTime    string `yaml:"max_conn_idle_time"`
		KeepAliveInterval  string `yaml:"keep_alive_interval"`
	} `yaml:"database"`
	Auth struct {
		PIN                   string `yaml:"pin"`
		AccessTokenSecretKey  string `yaml:"access_token_secret_key"`
		AccessTokenExpires    string `yaml:"access_token_expires"`
		RefreshTokenSecretKey string `yaml:"refresh_token_secret_key"`
		RefreshTokenExpires   string `yaml:"refresh_token_expires"`
		PermanentTokens       *bool  `yaml:"permanent_tokens"`
		BcryptCost            int    `yaml:"bcrypt_cost"`
		DefaultAdminEmail     string `yaml:"default_admin_email"`
		DefaultAdminPassword  string `yaml:"default_admin_password"`
	} `yaml:"auth"`
	Logging struct {
		Directory string `yaml:"directory"`
		File      string `yaml:"file"`
	} `yaml:"logging"`
	External struct {
		BaseURLCV       string `yaml:"base_url_cv"`
		ModelStorageDir string `yaml:"model_storage_dir"`
	} `yaml:"external"`
	Swagger struct {
		Enabled  *bool  `yaml:"enabled"`
		Host     string `yaml:"host"`
		BasePath string `yaml:"base_path"`
		Services string `yaml:"services"`
	} `yaml:"swagger"`
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

	setString := func(key, value string) {
		if value != "" {
			_ = os.Setenv(key, value)
		}
	}
	setBool := func(key string, value *bool) {
		if value != nil {
			_ = os.Setenv(key, fmt.Sprintf("%t", *value))
		}
	}
	setInt := func(key string, value int) {
		if value != 0 {
			_ = os.Setenv(key, fmt.Sprintf("%d", value))
		}
	}

	setString("APP_NAME", cfg.App.Name)
	setString("APP_ENV", cfg.App.Env)
	setString("APP_HOST", cfg.App.Host)
	setInt("APP_PORT", cfg.App.Port)
	setString("APP_URL", cfg.App.URL)
	setBool("APP_DEBUG", cfg.App.Debug)
	setBool("APP_PROMETHEUS", cfg.App.Prometheus)
	setString("APP_REQUEST_TIMEOUT", cfg.App.RequestTimeout)
	setString("APP_SHUTDOWN_TIMEOUT", cfg.App.ShutdownTimeout)
	setString("APP_SHUTDOWN_WAIT", cfg.App.ShutdownWait)
	setInt("APP_RATE_LIMIT_PER_SECOND", cfg.App.RateLimitPerSecond)
	setInt("APP_RATE_LIMIT_IP_PER_SECOND", cfg.App.RateLimitIPPerSecond)
	setString("APP_PUBLIC_API_KEY", cfg.App.PublicAPIKey)

	setString("DB_DRIVER", cfg.Database.Driver)
	setString("DB_SCHEMA", cfg.Database.Path)
	setString("DB_HOST", cfg.Database.Host)
	setString("DB_PORT", cfg.Database.Port)
	setString("DB_USERNAME", cfg.Database.Username)
	setString("DB_PASSWORD", cfg.Database.Password)
	setString("DB_SSLMODE", cfg.Database.SSLMode)
	setInt("DB_MAX_OPEN_CONNS", cfg.Database.MaxOpenConnections)
	setInt("DB_MAX_IDLE_CONNS", cfg.Database.MaxIdleConnections)
	setString("DB_MAX_LIFETIME_CONN", cfg.Database.MaxLifetime)
	setString("DB_MAX_CONN_IDLE_TIME", cfg.Database.MaxConnIdleTime)
	setString("DB_KEEP_ALIVE_INTERVAL_CONN", cfg.Database.KeepAliveInterval)

	setString("AUTH_PIN", cfg.Auth.PIN)
	setString("AUTH_ACCESS_TOKEN_SECRET_KEY", cfg.Auth.AccessTokenSecretKey)
	setString("AUTH_ACCESS_TOKEN_EXPIRES", cfg.Auth.AccessTokenExpires)
	setString("AUTH_REFRESH_TOKEN_SECRET_KEY", cfg.Auth.RefreshTokenSecretKey)
	setString("AUTH_REFRESH_TOKEN_EXPIRES", cfg.Auth.RefreshTokenExpires)
	setBool("AUTH_PERMANENT_TOKENS", cfg.Auth.PermanentTokens)
	setInt("AUTH_BCRYPT_COST", cfg.Auth.BcryptCost)
	setString("DEFAULT_ADMIN_EMAIL", cfg.Auth.DefaultAdminEmail)
	setString("DEFAULT_ADMIN_PASSWORD", cfg.Auth.DefaultAdminPassword)

	setString("LOG_DIR", cfg.Logging.Directory)
	setString("LOG_FILE", cfg.Logging.File)
	setString("BASE_URL_CV", cfg.External.BaseURLCV)
	setString("MODEL_STORAGE_DIR", cfg.External.ModelStorageDir)
	setBool("SWAGGER_ENABLED", cfg.Swagger.Enabled)
	setString("SWAGGER_HOST", cfg.Swagger.Host)
	setString("SWAGGER_BASE_PATH", cfg.Swagger.BasePath)
	setString("SWAGGER_SERVICES", cfg.Swagger.Services)

	return nil
}
