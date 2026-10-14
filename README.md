# Go Backend Template

Backend template in Go with a clean architecture structure inspired by the NATA HR service pattern.

## Structure

```text
cmd/
constant/
database/
  migrations/
  seeders/
docs/
internal/
  application/
    shared/dto/
    auth/
    menu/
    role/
    user/
    ...
  domain/
    contract/
      external/
      repository/
      usecase/
    model/
  infrastructure/
    api/
    database/
    logger/
    repository/
    runtime/
      container/
      env/
    validator/
  presentation/
    handler/
    middleware/
    presenter/
  test/
    mock/
      external/
      repository/
      usecase/
scripts/
```

## Environment

Copy `.env.example` to `.env` and adjust the values for your environment.

Important values:

- `APP_NAME`
- `APP_HOST`
- `APP_PORT`
- `DB_DRIVER`
- `DB_USERNAME`
- `DB_PASSWORD`
- `DB_SCHEMA`
- `DB_HOST`
- `DB_PORT`
- `SWAGGER_ENABLED`
- `SWAGGER_SERVICES`

## Run

```bash
go mod download
make run
```

## Database seed

```bash
make migrate.up
make migrate.seed
```

Dengan Docker Compose, PostgreSQL tersedia di host pada port `5437`. Jalankan `make docker.up` terlebih dahulu. Untuk menjalankan migration/seeder dari dalam container, gunakan `make docker.migrate.up` dan `make docker.migrate.seed`.

Dokumentasi handoff frontend untuk Master Source tersedia di [docs/frontend-handoff-master-source.md](docs/frontend-handoff-master-source.md).

## Swagger

Swagger UI is enabled with `SWAGGER_ENABLED=true`.
Use `SWAGGER_SERVICES` to filter the services shown in the docs.

## Test

```bash
go test ./...
```
