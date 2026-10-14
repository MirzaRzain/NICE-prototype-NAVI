package db

import (
	"database/sql"
	"log"
	"time"
)

func keepAlive(db *sql.DB, driver, schema string, interval time.Duration) {
	if interval <= 0 {
		return
	}
	for {
		err := db.Ping()
		if err != nil {
			log.Printf("ERROR db.keepAlive driver=%s schema=%s \n%s \n\n", driver, schema, err)
		}

		time.Sleep(interval)
	}
}

func startKeepAlive(db *sql.DB, driver, schema string, interval time.Duration) {
	if interval > 0 {
		go keepAlive(db, driver, schema, interval)
	}
}
