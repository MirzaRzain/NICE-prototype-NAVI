# Agent Guide — NAVI Backend

Panduan singkat untuk agent/developer yang bekerja di repository ini.

## Ringkasan

Backend REST API berbasis Go 1.24 dan Fiber v2, dengan struktur clean architecture sederhana:

- `cmd/`: entry point server dan seeder.
- `constant/`: error, permission, dan default configuration.
- `database/`: migration SQL dan seeder.
- `docs/`: generated Swagger.
- `internal/domain/`: model dan contract/interface.
- `internal/application/`: use case, DTO, mapper, JWT, utility.
- `internal/infrastructure/`: HTTP server, route, database, repository, runtime, logger.
- `internal/presentation/`: handler, presenter, middleware, response, validator.

Modul bisnis yang tersedia adalah authentication/session, user, role, menu, dan menu permission.

`RAW_EVENT_RECREATION_GUIDE.md` hanya blueprint alur raw-event camera service. Implementasi dan folder yang dirujuk dokumen itu belum terlihat pada tree repository ini.

## Arsitektur

Alur umum: HTTP route -> middleware -> presentation handler -> application use case -> domain contract -> infrastructure repository/query -> database.

Entry point server ada di `cmd/server/main.go`. Server me-load `.env`, membuat runtime context, logger, koneksi database, validator, container, lalu menjalankan Fiber melalui `internal/infrastructure/api/server.go`.

`internal/infrastructure/runtime/container.Container` membagikan database, logger, dan validator. Use case menerima `db.DB` lalu membuat repository dari `db.Query`, sehingga transaksi dan sqlmock dapat digunakan di test. Gunakan `internal/infrastructure/database.WithTransaction` untuk operasi multi-query yang harus atomik.

Aturan bisnis tetap di use case; SQL dan detail database tetap di infrastructure; handler tidak boleh memanggil SQL langsung.

## HTTP API

Base path adalah `/api/v1`; root health response ada di `/`. Swagger tersedia di `/docs/*` jika `SWAGGER_ENABLED=true`.

| Area | Endpoint |
| --- | --- |
| Auth | `POST /auth/login`, `POST /auth/logout`, `GET /auth/me` |
| Users | `GET/POST /users`, `GET/PUT/DELETE /users/:user_id`, `GET/PUT /users/me` |
| Roles | `GET/POST /roles`, `GET/PUT/DELETE /roles/:role_id`, `GET/PUT /roles/role-access/:role_id` |
| Menus | `GET/POST /menus`, `GET /menus/sidebar`, `GET /menus/parent`, `GET/PUT/DELETE /menus/:menu_id`, `PUT /menus/reorder` |
| Menu permissions | CRUD di `/menus/:menu_id/permissions` dan `/:menu_permission_id` |

Sebagian besar route selain login memakai `EnsureToken.ValidateToken()`. Middleware membaca `Authorization: Bearer <token>`, memverifikasi access token, lalu memastikan session masih ada di database. `Authorize(...)` memeriksa menu key dan permission code, termasuk inherited permission dari `constant/permission_hierarchy.go`.

Response sukses umumnya berbentuk `code`, `data`, dan `message`; endpoint list menambahkan `paginate`. Error terpusat di `internal/infrastructure/api/error_handler.go` dan memakai sentinel di `constant/constant_error.go`.

## Database

`internal/infrastructure/database` mendukung driver `postgresql`, `mysql`, dan `sqlite`. SQLite ditujukan untuk benchmark lokal; gunakan file database untuk benchmark concurrent write dan `:memory:` hanya untuk benchmark single-connection. Pool memakai `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, lifetime/idle timeout, dan SQLite diaktifkan dengan WAL serta busy timeout.

Migration master raw-event yang sudah tersedia mencakup `master_video`, `master_camera`, `master_camera_zone`, `master_usecase`, `usecase_zone`, `master_inference`, `master_object`, dan `usecase_inference_object`. Migration tersebut PostgreSQL-first agar konsisten dengan migration existing (`JSONB`, `BIGSERIAL`, dan `TIMESTAMPTZ`). SQLite memiliki migration terpisah satu tabel per pasangan file di `database/sqlite/migrations`, yang dijalankan dengan `make sqlite.migrate.up`, `make sqlite.migrate.down`, atau `make sqlite.setup`; jangan menjalankan `make migrate.up` ketika `DB_DRIVER=sqlite`.

Untuk SQLite di WSL, simpan file database di filesystem Linux (misalnya `/home/<user>/navi-data`) untuk menghindari masalah WAL/disk I/O pada `/mnt/c`. Driver akan fallback dari WAL ke journal mode `DELETE` jika filesystem tidak mendukung WAL.

Master API memakai workflow dua langkah: `POST /api/v1/master/cameras` mengembalikan `camera.id`, kemudian frontend memakai ID tersebut untuk preview dan CRUD zone melalui `GET/POST /api/v1/master/cameras/:camera_id/zones` serta `PUT/DELETE /api/v1/master/cameras/:camera_id/zones/:zone_id`. Request zone menerima `usecase_ids`, sehingga satu zone dapat di-assign ke banyak usecase dalam satu transaksi. Camera dapat diedit melalui `PUT /api/v1/master/cameras/:camera_id`; metadata yang didukung meliputi `area_location_id`, `usecase_config`, `channels`, `camera_type`, `capture_protocol`, dan `rtc_url`. Hanya `camera_code` yang unique untuk data aktif; field lain tidak diberi unique constraint. Video dibuat melalui `POST /api/v1/master/videos`.

Capture API memakai proxy `POST /api/v1/cv/capture` dengan alias `POST /api/v1/capture`. Proxy membaca `BASE_URL_CV`, meneruskan payload ke endpoint CV `/api/v1/capture`, menambahkan `channel` dari `master_camera.channels` bila belum ada, dan menghapus `source.camera.port`. Status code, content type, dan body dari CV diteruskan ke client.

Menu `Master Source` memakai endpoint gabungan `GET /api/v1/master/sources?page=1&limit=10&search=...` untuk pagination camera dan video. Item dibedakan melalui `source_type` (`camera` atau `video`); create camera dan create video tetap endpoint terpisah.

Master usecase memakai workflow dua tahap: `POST /api/v1/master/usecases` untuk data dasar, lalu `POST /api/v1/master/usecases/:usecase_id/inference-objects` untuk assign banyak pasangan inference-object. Edit data dasar memakai `PUT /api/v1/master/usecases/:usecase_id`; edit assignment memakai `PUT /api/v1/master/usecases/:usecase_id/inference-objects`, dengan daftar baru menggantikan assignment lama. Gunakan `GET` pada endpoint inference-objects untuk mengisi form edit.

Seeder didaftarkan di `database/seeders/init.go` untuk user, role, menu, user-role, menu permission, dan role access. Migration memakai pasangan `.up.sql` dan `.down.sql`; pertahankan timestamp prefix, foreign key, dan konvensi audit/soft-delete modul terkait.

## Environment

Konfigurasi aktual dibaca oleh `internal/infrastructure/runtime/config.go` dan `internal/infrastructure/database/database_option.go`.

Minimal database environment:

```text
DB_DRIVER       postgresql, mysql, atau sqlite
DB_USERNAME
DB_PASSWORD
DB_SCHEMA       nama database/schema; untuk sqlite isi path file atau `:memory:`
DB_HOST
DB_PORT
DB_SSLMODE
```

Environment aplikasi penting: `APP_NAME`, `APP_HOST`, `APP_PORT`, `APP_ENV`, `APP_DEBUG`, `APP_REQUEST_TIMEOUT`, `APP_SHUTDOWN_TIMEOUT`, `APP_SHUTDOWN_WAIT`, `APP_RATE_LIMIT_PER_SECOND`, `APP_RATE_LIMIT_IP_PER_SECOND`, `AUTH_ACCESS_TOKEN_SECRET_KEY`, `AUTH_REFRESH_TOKEN_SECRET_KEY`, `BASE_URL_CV`, `SWAGGER_ENABLED`, `SWAGGER_HOST`, `SWAGGER_BASE_PATH`, dan `SWAGGER_SERVICES`.

Logging aplikasi memakai `LOG_DIR` dan `LOG_FILE`, default ke `logs/app.txt`. Logger menulis JSON terstruktur ke stdout dan file secara append. HTTP request/response dan error handler dicatat oleh middleware; error service/repository dicatat dari layer use case/repository melalui logger yang sama. Jangan mencatat password, access token, refresh token, atau request body sensitif.

Untuk edge device dengan satu akun, tetap pertahankan `app_user`, `app_login`, dan `sessions` karena ketiganya menyimpan identitas, audit login, serta lifecycle session. Kredensial superadmin sebaiknya hanya dipakai sebagai bootstrap pertama dari environment/device provisioning, kemudian password di-hash ke `app_user`; jangan menjadikan env sebagai source of truth untuk login setiap request.

Catatan: README menyebut `DB_NAME`, tetapi kode aktual membutuhkan `DB_SCHEMA`. Untuk connection pool, gunakan nama plural `DB_MAX_IDLE_CONNS` dan `DB_MAX_OPEN_CONNS` seperti pada `.env.example`. `DB_HOST`, `DB_PORT`, username, password tidak diperlukan untuk SQLite.

## Menjalankan

Pastikan `.env` tersedia. Perintah langsung:

```bash
go mod download
go run ./cmd/server
go run ./cmd/seeder --gseed
go test ./...
```

`Makefile` aktif menyediakan `make migrate.up` dan `make migrate.seed` untuk eksekusi dari host/WSL. Dengan Docker Compose, PostgreSQL dipetakan ke host port `5437`, sehingga target host memakai `DB_MIGRATE_PORT=5437` secara default. Target `make docker.migrate.up` dan `make docker.migrate.seed` menjalankan migration/seeder di dalam container dan memakai port internal PostgreSQL `5432`. Override port host untuk PostgreSQL lokal dengan `make migrate.up DB_MIGRATE_PORT=5432`.

Docker Compose menyediakan app dan PostgreSQL. PostgreSQL memakai port lokal `5437`; dari container app host database adalah `postgres` dengan port `5432`.

## Testing dan perubahan kode

Test dominan berada di `internal/application/**/usecase/*_test.go`, memakai `testify`, `sqlmock`, dan mock repository di `internal/test/mock/repository`.

Saat menambah fitur:

1. Tambahkan/update domain model dan contract bila ada boundary baru.
2. Tambahkan DTO/mapper dan use case untuk aturan bisnis.
3. Tambahkan query/repository implementation serta mock contract.
4. Tambahkan presenter request/response dan handler.
5. Daftarkan route dan middleware/authorization.
6. Tambahkan migration/seeder hanya bila schema/data awal diperlukan.
7. Tambahkan unit test untuk happy path, business error, dan database error.
8. Jalankan `gofmt` dan `go test ./...`.

Gunakan context propagation dan error sentinel yang sudah ada agar status HTTP tetap konsisten.

### Aturan migration

- Satu tabel = satu migration file `.up.sql` dan satu pasangan `.down.sql` dengan timestamp prefix yang sama.
- Jika mengubah atau melakukan `ALTER TABLE`, jangan mengedit migration yang sudah pernah dibuat/berjalan. Buat migration baru khusus untuk perubahan tabel tersebut, tetap satu tabel per migration.
- Index, constraint, dan trigger yang menjadi bagian dari pembuatan tabel boleh berada di migration tabel itu.
- Migration baru harus mengikuti urutan dependency foreign key dan menyediakan rollback yang aman.

## Hal penting

- Jangan commit secret JWT atau kredensial.
- Token valid tetap harus punya session aktif di database.
- Gunakan konstanta menu/action untuk authorization jika tersedia.
- Beberapa route update/order saat ini tidak memiliki `Authorize(...)` langsung di file route. Pertahankan perilaku existing kecuali task meminta perubahan authorization; tambahkan regression test bila mengubahnya.
- `docs/swagger.*` adalah generated files; regenerate setelah mengubah anotasi endpoint dan review diff.
- Ada typo historis seperti `update_menu_order_hanlders.go`; cari berdasarkan nama file yang benar-benar ada.
- Cek diff/status sebelum mengubah file yang sudah dimodifikasi user. Environment agent tidak selalu menyediakan working tree Git.

## File rujukan cepat

- Bootstrap: `cmd/server/main.go`, `cmd/seeder/main.go`
- Routes: `internal/infrastructure/api/routes/routes.go`
- Server: `internal/infrastructure/api/server.go`
- Config: `internal/infrastructure/runtime/config.go`
- Auth middleware: `internal/presentation/middleware/auth_middleware.go`
- Error mapping: `internal/infrastructure/api/error_handler.go`
- Transaction: `internal/infrastructure/database/trx.go`
- Contracts: `internal/domain/contract`
- Use cases: `internal/application`
- Repositories: `internal/infrastructure/repository`
- Schema: `database/migrations`
- Seeds: `database/seeders`
