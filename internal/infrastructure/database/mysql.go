package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(opt *databaseOption) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		opt.username,
		opt.password,
		opt.host,
		opt.port,
		opt.schema,
	)

	db, err = openSQL("mysql", dsn, opt.connectionOption)
	if err != nil {
		return
	}

	log.Printf("Successfully connected to MySQL %s:%d database: %s", opt.host, opt.port, opt.schema)

	startKeepAlive(db, opt.driver, opt.schema, opt.keepAliveInterval)

	return
}
