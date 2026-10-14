package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "modernc.org/sqlite"
)

// NewSQLite opens a SQLite database suitable for local throughput benchmarks.
func NewSQLite(opt *databaseOption) (db *sql.DB, err error) {
	db, err = openSQL("sqlite", sqliteDSN(opt.schema), opt.connectionOption)
	if err != nil {
		return
	}

	if _, err = db.ExecContext(opt.ctx, `PRAGMA foreign_keys=ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err = db.ExecContext(opt.ctx, `PRAGMA journal_mode=WAL`); err != nil {
		log.Printf("SQLite WAL unavailable for %s, falling back to DELETE journal mode: %v", opt.schema, err)
		if _, fallbackErr := db.ExecContext(opt.ctx, `PRAGMA journal_mode=DELETE`); fallbackErr != nil {
			_ = db.Close()
			return nil, fallbackErr
		}
	}
	if _, err = db.ExecContext(opt.ctx, `PRAGMA synchronous=NORMAL`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err = db.ExecContext(opt.ctx, `PRAGMA busy_timeout=5000`); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Printf("Successfully connected to sqlite database: %s", opt.schema)
	return
}

func sqliteDSN(path string) string {
	if path == ":memory:" {
		return "file:navi-benchmark?mode=memory&cache=shared&_pragma=busy_timeout(5000)"
	}
	if strings.HasPrefix(path, "file:") {
		return path
	}
	return fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", url.PathEscape(path))
}
