#!/bin/sh
set -eu

cmd="${1:-server}"

shift || true

db_url() {
  host="${DB_HOST:-postgres}"
  port="${DB_PORT:-5432}"
  user="${DB_USERNAME:-postgres}"
  pass="${DB_PASSWORD:-postgres}"
  name="${DB_SCHEMA:-go-backend-template}"

  printf "postgres://%s:%s@%s:%s/%s?sslmode=disable" "$user" "$pass" "$host" "$port" "$name"
}

case "$cmd" in
  server)
    exec /app/go-backend-template
    ;;

  migrate.up)
    exec /usr/local/bin/migrate -path /app/database/migrations -database "$(db_url)" up
    ;;

  migrate.down)
    if [ "${1:-}" = "" ]; then
      exec /usr/local/bin/migrate -path /app/database/migrations -database "$(db_url)" down
    fi

    exec /usr/local/bin/migrate -path /app/database/migrations -database "$(db_url)" down "$1"
    ;;

  migrate.fix)
    if [ "${1:-}" = "" ]; then
      echo "usage: migrate.fix <version>" >&2
      exit 2
    fi

    exec /usr/local/bin/migrate -path /app/database/migrations -database "$(db_url)" force "$1"
    ;;

  migrate.seed)
    exec /app/go-backend-template-seed --gseed
    ;;

  sqlite.schema)
    echo "SQLite schema command is intended for the host development workflow." >&2
    exit 2
    ;;

  *)
    exec "$cmd" "$@"
    ;;
esac
