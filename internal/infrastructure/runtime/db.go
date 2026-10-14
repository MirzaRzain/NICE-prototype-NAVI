package config

import "time"

type DBConfig struct {
	Driver            string
	Username          string
	Password          string
	Name              string
	Host              string
	Port              string
	SSLMode           string
	MaxOpenConn       int
	MaxIdleConn       int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	KeepAliveInterval time.Duration
}
