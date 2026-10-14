package constant

import "time"

// middleware.
const (
	DefaultMdwHeaderToken         = "Authorization"
	DefaultMdwHeaderBearer        = "Bearer"
	DefaultMdwRateLimiter         = 20
	DefaultMdwRateLimiterDuration = time.Minute
	DefaultMdwTimeout             = 10 * time.Second
)

// runtime.
const (
	DefaultAppName            = "GO Backend Template API"
	DefaultAppPort            = 8000
	DefaultAppHost            = "127.0.0.1"
	DefaultAppShutdownTimeout = 10 * time.Second
	DefaultAppShutdownWait    = 200 * time.Millisecond
)

// db connection.
const (
	DefaultDBSSLMode           = "disable"
	DefaultDBMaxOpenConns      = 30
	DefaultDBMaxIdleConns      = 30
	DefaultDBMaxIdleTime       = 15 * time.Second
	DefaultDBMaxLifetime       = 60 * time.Second
	DefaultDBMaxConnWaitTime   = 1 * time.Second
	DefaultDBMaxConnLifetime   = 3600 * time.Second
	DefaultDBMaxConnIdleTime   = 60 * time.Second
	DefaultDBKeepAliveInterval = 3 * time.Minute
)

// pagination.
const (
	DefaultOrder = "created_at DESC"
	DefaultPage  = 1
	DefaultLimit = 10
)
