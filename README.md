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

Copy `.env.example` to `.env` and adjust the values for your environment. For SSH development without `.env`, configure `config.yaml`; when that file exists, it is loaded as the primary configuration source.

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

## Deploy dengan Docker Compose di server SSH

Di server, siapkan environment Docker lalu jalankan:

```bash
cp .env.docker.example .env
# edit .env: password database, secret JWT, APP_URL, dan BASE_URL_CV
docker compose config --quiet
docker compose up -d --build
docker compose run --rm app migrate.up
docker compose run --rm app migrate.seed
docker compose ps
docker compose logs -f app
```

Container aplikasi bind ke `0.0.0.0` dan PostgreSQL hanya diperlukan oleh container internal. Port API mengikuti `APP_PORT` (default `8000`); buka port tersebut di firewall/security group bila API perlu diakses dari luar server.

`config.yaml` adalah konfigurasi SQLite untuk development lokal dan sengaja tidak dipakai di image Docker. Gunakan `.env.docker.example` untuk deployment Compose.

Dokumentasi handoff frontend untuk Master Source tersedia di [docs/frontend-handoff-master-source.md](docs/frontend-handoff-master-source.md).

Dokumentasi handoff API untuk integrasi frontend tersedia di [docs/api-handoff.md](docs/api-handoff.md).

## Swagger

Swagger UI is enabled with `SWAGGER_ENABLED=true`.
Use `SWAGGER_SERVICES` to filter the services shown in the docs.

## Test

```bash
go test ./...
```
