package config

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/infrastructure/runtime/env"
)

type Swagger struct {
	Host      string
	BasePath  string
	Services  []string
	IsEnabled bool
}

type Config struct {
	Name                    string
	Host                    string
	Port                    int
	URL                     string
	Debug                   bool
	Prometheus              bool
	RequestTimeout          time.Duration
	ShutdownTimeoutDuration time.Duration
	ShutdownWaitDuration    time.Duration
	RateLimitPerSecond      int
	RateLimitIPPerSecond    int
	PublicAPIKey            string
	BaseURLCV               string
	DB                      DBConfig
	Swagger                 Swagger
}

func NewConfig() *Config {
	// DB Config load
	DBConfig := new(DBConfig)
	DBConfig.Driver = env.MustGetEnv("DB_DRIVER")
	DBConfig.Username = env.GetString("DB_USERNAME", "")
	DBConfig.Password = env.GetString("DB_PASSWORD", "")
	DBConfig.Name = env.GetString("DB_SCHEMA", "")
	DBConfig.Host = env.GetString("DB_HOST", "")
	DBConfig.Port = env.GetString("DB_PORT", "")
	if DBConfig.Driver != "sqlite" && DBConfig.Name == "" {
		DBConfig.Name = env.MustGetEnv("DB_SCHEMA")
	}
	if DBConfig.Driver != "sqlite" && DBConfig.Host == "" {
		DBConfig.Host = env.MustGetEnv("DB_HOST")
	}
	DBConfig.SSLMode = env.GetString("DB_SSLMODE", constant.DefaultDBSSLMode)
	DBConfig.MaxOpenConn = env.GetInt("DB_MAX_OPEN_CONNS", constant.DefaultDBMaxOpenConns)
	DBConfig.MaxIdleConn = env.GetInt("DB_MAX_IDLE_CONNS", constant.DefaultDBMaxIdleConns)
	DBConfig.MaxConnLifetime = env.GetDuration("DB_MAX_LIFETIME_CONN", constant.DefaultDBMaxConnLifetime)
	DBConfig.MaxConnIdleTime = env.GetDuration("DB_MAX_CONN_IDLE_TIME", constant.DefaultDBMaxConnIdleTime)
	DBConfig.KeepAliveInterval = env.GetDuration("DB_KEEP_ALIVE_INTERVAL_CONN", constant.DefaultDBKeepAliveInterval)

	// swagger
	services := env.GetString("SWAGGER_SERVICES", "")
	swagger := Swagger{
		Host:      env.GetString("SWAGGER_HOST", ""),
		BasePath:  env.GetString("SWAGGER_BASE_PATH", ""),
		Services:  parseSwaggerServices(services),
		IsEnabled: env.GetBool("SWAGGER_ENABLED", false),
	}

	// Application config load
	requestTimeout := env.GetDuration("APP_REQUEST_TIMEOUT", constant.DefaultMdwTimeout)
	cfg := new(Config)
	cfg.Name = env.GetString("APP_NAME", constant.DefaultAppName)
	cfg.Host = env.GetString("APP_HOST", constant.DefaultAppHost)
	cfg.Port = env.GetInt("APP_PORT", constant.DefaultAppPort)
	cfg.URL = env.GetString("APP_URL", "")
	cfg.Debug = env.GetBool("APP_DEBUG", false)
	cfg.Prometheus = env.GetBool("APP_PROMETHEUS", false)
	cfg.RequestTimeout = requestTimeout
	cfg.ShutdownTimeoutDuration = env.GetDuration("APP_SHUTDOWN_TIMEOUT", constant.DefaultAppShutdownTimeout)
	cfg.ShutdownWaitDuration = env.GetDuration("APP_SHUTDOWN_WAIT", constant.DefaultAppShutdownWait)
	cfg.RateLimitPerSecond = env.GetInt("APP_RATE_LIMIT_PER_SECOND", constant.DefaultMdwRateLimiter)
	cfg.RateLimitIPPerSecond = env.GetInt("APP_RATE_LIMIT_IP_PER_SECOND", constant.DefaultMdwRateLimiter)
	cfg.PublicAPIKey = env.GetString("APP_PUBLIC_API_KEY", "")
	cfg.BaseURLCV = env.GetString("BASE_URL_CV", "")
	cfg.DB = *DBConfig
	cfg.Swagger = swagger

	return cfg
}

func (c *Config) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func parseSwaggerServices(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	services := make([]string, 0, len(parts))
	seen := map[string]struct{}{}

	for _, part := range parts {
		service := strings.TrimSpace(part)
		if service == "" {
			continue
		}

		if _, exists := seen[service]; exists {
			continue
		}

		seen[service] = struct{}{}
		services = append(services, service)
	}

	return services
}

func LoadConfig() error {
	if _, err := os.Stat("config.yaml"); err == nil {
		return loadConfigFile("config.yaml")
	}
	return godotenv.Load(".env")
}

func NewRuntimeContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}
