# Device deployment

Backend device menggunakan SQLite dan konfigurasi YAML. Salin `config.yaml.example` menjadi `config.yaml`, lalu ubah PIN dan secret sebelum deploy.

```yaml
database:
  driver: sqlite
  path: /home/mirza-ruzain/navi-data/navi-benchmark.db

auth:
  pin: "123456"
  access_token_secret_key: change-this-access-secret
  refresh_token_secret_key: change-this-refresh-secret
  permanent_tokens: true
```

Login menggunakan satu PIN global:

```http
POST /api/v1/auth/login
Content-Type: application/json

{"pin":"123456"}
```

Access token dan refresh token tidak memiliki expiry JWT saat `permanent_tokens: true`. Endpoint menu, role, permission, dan user tidak diregistrasikan pada device build.

Jalankan schema dan seed SQLite dengan:

```bash
make sqlite.setup
```
