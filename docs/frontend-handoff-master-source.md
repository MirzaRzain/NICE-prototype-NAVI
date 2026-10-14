# Frontend Handoff — Master Source API

Base URL local: `http://127.0.0.1:8008/api/v1`

Semua endpoint selain login memakai:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

## Response

Object sukses:

```json
{"code":200,"data":{},"message":"Success"}
```

Pagination:

```json
{
  "code": 200,
  "data": [],
  "message": "Success",
  "paginate": {"current_page":1,"per_page":10,"total_page":1,"total_data":0}
}
```

Error:

```json
{"code":400,"message":"Invalid request"}
```

## Login Device (PIN)

```http
POST /auth/login
```

```json
{"pin":"<PIN dari config.yaml>"}
```

Simpan `data.access_token` sebagai Bearer token:

```http
Authorization: Bearer <access_token>
```

Response login tidak lagi memiliki user karena device tidak memakai `app_user`:

```json
{
  "code": 200,
  "data": {
    "access_token": "<jwt permanen>",
    "refresh_token": "<jwt permanen>"
  },
  "message": "Successfully login"
}
```

Access token dan refresh token bersifat permanen pada deployment device. Frontend tidak perlu mengandalkan expiry untuk memaksa login ulang. Jika API mengembalikan `401`, hapus token lokal dan arahkan kembali ke form PIN.

Endpoint `/auth/me` tidak tersedia.

## Endpoint yang Dihapus pada Device

Frontend tidak perlu membuat halaman atau request untuk:

- `/users`
- `/menus`
- `/roles`
- menu permission dan role access

Tidak ada lagi user, role, menu, permission, atau sidebar berbasis `app_menu`.

## Master Source gabungan

### GET `/master/sources?page=1&limit=10&search=camera&source_type=camera&start_time=2026-10-08T00:00:00Z&end_time=2026-10-09T00:00:00Z`

`search` mencari pada `code` dan `name`. Gunakan `source_type` untuk membedakan camera dan video.
`source_type` bersifat opsional dan hanya menerima `camera` atau `video`; jika tidak dikirim, hasil berasal dari kedua tabel.
`start_time` dan `end_time` bersifat opsional, memfilter `created_at`. Keduanya harus RFC3339; frontend disarankan mengirim UTC dengan suffix `Z`. `start_time` inklusif, sedangkan `end_time` eksklusif.

```json
{
  "code": 200,
  "data": [
    {"id":1,"source_type":"camera","code":"CAM-001","name":"Lobby Camera","type":"camera","camera_status":true,"ip_address":"192.168.1.10","stream_url":"rtsp://192.168.1.10/live","created_at":"2026-10-08T10:00:00Z"},
    {"id":2,"source_type":"video","code":"VID-001","name":"Demo Video","type":"mp4","source_video":"/videos/input.mp4","result_video":null,"created_at":"2026-10-08T10:05:00Z"}
  ],
  "message":"Successfully get master sources",
  "paginate":{"current_page":1,"per_page":10,"total_page":1,"total_data":2}
}
```

## Master Video

### POST `/master/videos`

```json
{
  "video_code":"VID-001",
  "video_name":"Demo Video",
  "video_type":"mp4",
  "source_video":"/videos/input.mp4",
  "result_video":null
}
```

`result_video` nullable. Response `201` mengembalikan `data.id` dan object video.

## Master Camera

### POST `/master/cameras`

```json
{
  "area_location_id": 1,
  "camera_code":"CAM-001",
  "camera_name":"Lobby Camera",
  "camera_status":true,
  "ip_address":"192.168.1.10",
  "stream_url":"rtsp://192.168.1.10/live",
  "auth_type":"basic",
  "username":"camera-user",
  "password":"camera-password",
  "installed_at":"2026-10-08T10:00:00Z",
  "usecase_config":[],
  "channels":"channel-1",
  "camera_type":"fixed",
  "capture_protocol":"rtsp",
  "capture_mode":"STREAM",
  "rtc_url":"https://example.local/camera/CAM-001",
  "device_path":null
}
```

Semua field input camera opsional untuk sementara. `camera_code` satu-satunya business field yang memiliki unique constraint jika diisi; nilai `NULL` diperbolehkan. Field lain tidak unique dan dapat `null`, termasuk `area_location_id`, `ip_address`, `stream_url`, dan `device_path`. `id` tetap primary key dan dibuat database.

Untuk USB camera, frontend dapat mengirim:

```json
{
  "camera_type":"usb",
  "device_path":"0"
}
```

`ip_address` dan `stream_url` tidak perlu dikirim untuk USB. Untuk camera RTSP/network, isi `ip_address`, `stream_url`, dan `capture_protocol` bila tersedia. `capture_mode` bersifat opsional dan hanya menggunakan nilai `STREAM` atau `SNAPSHOT`. `installed_at` menggunakan RFC3339, `usecase_config` default `[]`, dan `password` tidak pernah dikembalikan pada response.

`capture_mode` dipisahkan dari `camera_type` dan `capture_protocol`: `camera_type` menjelaskan perangkat (misalnya `usb`), sedangkan `capture_protocol` menjelaskan protokol koneksi (misalnya `rtsp` atau `onvif`).

### PUT `/master/cameras/:camera_id`

Body sama dengan create camera. Endpoint ini memperbarui metadata camera dan source fields. `password` boleh tidak dikirim; jika tidak dikirim, password lama tetap digunakan. Response `200`:

```json
{
  "code":200,
  "data":null,
  "message":"Successfully update master camera"
}
```

`camera_code` memiliki unique constraint untuk camera yang belum dihapus. Tidak ada unique constraint lain pada field camera. `camera_type` dapat berupa `usb` atau tipe lain yang didukung database/backend.

## Camera Zone

### Alur stepper

1. `POST /master/cameras`.
2. Simpan `response.data.id` sebagai `camera_id`.
3. Preview: `GET /master/cameras/:camera_id/zones`.
4. Tambah: `POST /master/cameras/:camera_id/zones`.
5. Edit: `PUT /master/cameras/:camera_id/zones/:zone_id`.
6. Hapus: `DELETE /master/cameras/:camera_id/zones/:zone_id`.

### POST `/master/cameras/:camera_id/zones`

```json
{
  "zone_name":"Lobby",
  "zone_type":"polygon",
  "zone_coordinate":{"points":[[0.1,0.1],[0.9,0.1],[0.9,0.9],[0.1,0.9]]},
  "config":{"max_person_count":6},
  "usecase_ids":[1,2,3]
}
```

`usecase_ids` mendukung multiple assignment dan diproses dalam satu transaksi.

### GET `/master/cameras/:camera_id/zones`

`data` berupa array object zone:

```json
{
  "id":10,
  "camera_id":1,
  "zone_name":"Lobby",
  "zone_type":"polygon",
  "zone_coordinate":{"points":[[0.1,0.1],[0.9,0.1],[0.9,0.9],[0.1,0.9]]},
  "config":{"max_person_count":6},
  "usecases":[{"id":1,"name":"Zone Density","code":"ZONE_DENSITY"}],
  "created_at":"2026-10-08T10:10:00Z"
}
```

Body PUT sama seperti POST. `usecase_ids` menggantikan assignment lama. Response PUT/DELETE menggunakan `data:null` dan message sukses.

## Capture API

Camera service menyediakan proxy ke CV service. Frontend dapat memakai salah satu endpoint berikut:

```http
POST /api/v1/cv/capture
POST /api/v1/capture
```

Header:

```http
Content-Type: application/json
Accept: application/json
```

Contoh request:

```json
{
  "session_id": 9,
  "usecase": {
    "id": 10,
    "code": "people_monitoring",
    "config": {"minimum_duration": 3}
  },
  "source": {
    "type": "stream",
    "fps": 30,
    "camera": {
      "camera_id": "CAM-001",
      "stream_address": "rtsp://192.168.100.85:554/stream",
      "ip_address": "192.168.100.85",
      "capture_protocol": "onvif",
      "port": 554,
      "auth": {"type": "basic", "username": "admin", "password": "admin123"}
    }
  },
  "inference": {
    "type": "detection",
    "objects": ["person"],
    "config": {"confidence_threshold": 0.5, "iou_threshold": 0.45}
  },
  "zones": [
    {
      "zone_id": 12,
      "type": "polygon",
      "points": [[0.1,0.2],[0.8,0.2],[0.8,0.9],[0.1,0.9]],
      "config": {"minimum_duration": 3}
    }
  ],
  "output": {"mode": "event", "interval": 1}
}
```

Backend akan meneruskan request ke `${BASE_URL_CV}/capture` jika `BASE_URL_CV` sudah berakhiran `/api/v1`; jika belum, backend memakai `${BASE_URL_CV}/api/v1/capture`.

Jika `source.camera` tersedia atau `camera_id` tersedia di root dan `channel` belum dikirim, backend mencari `master_camera.channels` berdasarkan camera ID atau camera code lalu menambahkan `channel`. Field `source.camera.port` selalu dihapus sebelum request diteruskan. Response dari CV diteruskan apa adanya, termasuk status code, content type, dan body.

Error dari proxy:

- `400`: JSON capture tidak valid.
- `500`: `BASE_URL_CV` belum dikonfigurasi.
- `502`: gagal menghubungi CV service atau gagal membaca response-nya.

## Master Usecase

Inference dan object harus sudah tersedia di `master_inference` dan `master_object`.

Seed default:

| id | inference | code |
|---:|---|---|
| 1 | Detection | `detection` |
| 2 | Tracking | `tracking` |
| 3 | Recognition | `recognition` |

Object default: `Face` (`face`) dan `Person` (`person`).

Flow create/edit menggunakan dua data terpisah:

1. `POST /master/usecases` membuat atau mengubah data dasar di `master_usecase`.
2. `POST /master/usecases/:usecase_id/inference-objects` membuat assignment awal di `usecase_inference_object`.
3. `PUT /master/usecases/:usecase_id` mengedit data dasar `master_usecase`.
4. `PUT /master/usecases/:usecase_id/inference-objects` mengganti assignment `usecase_inference_object`.

Create usecase tidak otomatis membuat row di `usecase_inference_object`; assignment dilakukan pada step kedua.

### Step 1 — POST `/master/usecases`

```json
{"usecase_name":"Zone Density","usecase_code":"ZONE_DENSITY","description":"Menghitung kepadatan orang per zone"}
```

Response `201` mengembalikan `data.id` sebagai `usecase_id`.

### Step 2 — POST `/master/usecases/:usecase_id/inference-objects`

```json
{
  "inference_objects":[
    {"inference_id":1,"object_id":1},
    {"inference_id":2,"object_id":3}
  ]
}
```

Response `200`:

```json
{"code":200,"data":null,"message":"Successfully set inference objects"}
```

### Edit data dasar — PUT `/master/usecases/:usecase_id`

Body sama seperti POST create. Response `200` dengan `data:null`.

### Preview assignment — GET `/master/usecases/:usecase_id/inference-objects`

```json
{
  "code":200,
  "data":[{"id":10,"inference_id":1,"object_id":1,"inference_name":"People Detection","inference_code":"PEOPLE_DETECTION","object_name":"Person","object_code":"PERSON"}],
  "message":"Successfully get inference objects"
}
```

### Edit assignment — PUT `/master/usecases/:usecase_id/inference-objects`

Body sama seperti POST step 2. Daftar baru menggantikan seluruh assignment lama. Untuk menghapus semua:

```json
{"inference_objects":[]}
```

## Error

| HTTP | Arti |
|---:|---|
| 400 | JSON, parameter, atau validasi salah |
| 401 | Token tidak ada/expired/session invalid |
| 404 | Resource tidak ditemukan |
| 500 | Error database/server |

## Local SQLite

Konfigurasi device dibaca dari `config.yaml`, bukan `.env`:

```yaml
database:
  driver: sqlite
  path: /home/mirza-ruzain/navi-data/navi-benchmark.db

auth:
  pin: "123456"
  permanent_tokens: true
```

```bash
mkdir -p ~/navi-data
make sqlite.migrate.up
make migrate.seed
make run
```

Untuk deployment baru, gunakan `make sqlite.setup` agar migration dan seed dijalankan. Migration PostgreSQL tidak digunakan lagi.

Rollback migration terakhir: `make sqlite.migrate.down`.

Health check:

```bash
curl http://127.0.0.1:8000/
```
