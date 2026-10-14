package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresql(opt *databaseOption) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		opt.host,
		opt.port,
		opt.username,
		opt.password,
		opt.schema,
		opt.sslmode,
	)

	db, err = openSQL("postgres", dsn, opt.connectionOption)
	if err != nil {
		return
	}

	log.Printf("Successfully connected to postgresql %s:%d schema: %s", opt.host, opt.port, opt.schema)

	startKeepAlive(db, opt.driver, opt.schema, opt.keepAliveInterval)

	return
}
