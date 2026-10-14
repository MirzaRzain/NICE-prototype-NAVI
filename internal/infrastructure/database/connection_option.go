package db

import (
	"time"

	"be-dashboard-nba/constant"
)

type connectionOption struct {
	maxIdle           int
	maxOpen           int
	maxLifetime       time.Duration
	maxConnIdleTime   time.Duration
	keepAliveInterval time.Duration
}

func defaultConnectionOption() *connectionOption {
	return &connectionOption{
		maxIdle:           constant.DefaultDBMaxIdleConns,
		maxOpen:           constant.DefaultDBMaxOpenConns,
		maxLifetime:       constant.DefaultDBMaxConnLifetime,
		maxConnIdleTime:   constant.DefaultDBMaxConnIdleTime,
		keepAliveInterval: constant.DefaultDBKeepAliveInterval,
	}
}
