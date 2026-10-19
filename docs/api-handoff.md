# API Handoff — Frontend

Dokumen handoff untuk integrasi frontend dengan backend NICE Prototype Navi.

> Status: sinkron dengan route yang saat ini terdaftar di backend. Endpoint yang handler-nya sudah ada tetapi belum diregistrasikan ditandai jelas.

## Base URL

```text
http://127.0.0.1:8008/api/v1
```

Sesuaikan port dengan `APP_PORT` atau deployment.

## Header

Untuk request JSON:

```http
Content-Type: application/json
Accept: application/json
```

Untuk endpoint yang memang memakai session:

```http
Authorization: Bearer <access_token>
```

## Format response

Sukses:

```json
{"code":200,"data":{},"message":"Success"}
```

Paginated:

```json
{
  "code": 200,
  "data": [],
  "message": "Success",
  "paginate": {"current_page": 1, "per_page": 10, "total_page": 1, "total_data": 0}
}
```

Error:

```json
{"code":400,"message":"Invalid request","error":{}}
```

Field `error` bersifat opsional.

## Authentication dan kebutuhan token

Kondisi route runtime saat ini:

| Endpoint | Bearer token |
|---|---|
| `POST /auth/login` | Tidak perlu |
| `POST /auth/logout` | Wajib |
| `/master/*` yang saat ini terdaftar | Tidak diwajibkan oleh route |
| `POST /capture` dan `POST /cv/capture` | Tidak perlu |

Jadi frontend tidak perlu mengirim Bearer token untuk login, endpoint master yang aktif, atau capture. Token tetap disimpan karena dibutuhkan untuk logout dan dapat menjadi dasar pengamanan route pada deployment berikutnya.

### POST `/auth/login`

```json
{"pin":"<PIN dari config.yaml>"}
```

Simpan `data.access_token` untuk request berikutnya. Jika menerima `401`, hapus token lokal dan arahkan user kembali ke form PIN.

### POST `/auth/logout`

Membatalkan session aktif dan membutuhkan Bearer token. Response sukses menggunakan `data: null`.

> `/auth/me` dan login email/password masih tercantum di sebagian Swagger lama, tetapi tidak diregistrasikan oleh route device saat ini.

## Master Usecase

## Tiga API utama frontend

Kontrak utama frontend untuk pemilihan dan pengambilan model:

1. `GET /master/usecases/:usecase_code/models`
   Backend mencari usecase berdasarkan `usecase_code`, mengambil seluruh `inference_object` yang terhubung ke usecase tersebut, lalu mengembalikan semua model yang terhubung ke inference-object tersebut melalui `model_inference_object`.
2. `GET /master/models/download/:usecase_code`
   Backend memakai resolusi usecase/inference-object yang sama dan mengunduh model terbaru yang tersedia untuk usecase tersebut.
3. `GET /master/usecases`
   Mengembalikan daftar aktif dengan field `usecase_id`, `usecase_name`, dan `usecase_code`.

Endpoint download berdasarkan `model_id` dan `inference_object_id` tetap tersedia untuk kompatibilitas, tetapi bukan alur utama frontend.

### GET `/master/usecases`

Daftar usecase aktif:

```json
{
  "code": 200,
  "data": [{"usecase_id": 1, "usecase_code": "ZONE_DENSITY", "usecase_name": "Zone Density"}],
  "message": "Successfully get master usecases"
}
```

### GET `/master/usecases/config`

Mengembalikan usecase dan konfigurasi key-value:

```json
{
  "code": 200,
  "data": [{"usecase_id": 1, "usecase_code": "ZONE_DENSITY", "usecase_name": "Zone Density", "usecase_config": [{"key": "minimum_duration", "value": "3"}]}],
  "message": "Successfully get usecases with config"
}
```

### GET `/master/usecases/:usecase_id`

Mengambil detail usecase.

### POST `/master/usecases`

```json
{
  "usecase_name": "Zone Density",
  "usecase_code": "ZONE_DENSITY",
  "description": "Menghitung kepadatan orang per zone"
}
```

`usecase_name` dan `usecase_code` wajib, maksimal 255 karakter. Response `201` mengembalikan object usecase pada `data`, termasuk `data.id`.

### PUT `/master/usecases/:usecase_id`

Body sama dengan create. Response `200` mengembalikan `data: null`.

### GET `/master/usecases/:usecase_id/inference-objects`

Preview assignment inference-object. Item berisi `id` sebagai mapping usecase, serta `inference_object_id` sebagai ID kombinasi global yang dipakai saat assign model.

### POST/PUT `/master/usecases/:usecase_id/inference-objects`

POST membuat assignment awal; PUT mengganti seluruh assignment lama.

```json
{
  "inference_objects": [
    {"inference_id": 1, "object_id": 1},
    {"inference_id": 2, "object_id": 3}
  ]
}
```

Untuk menghapus semua assignment, kirim `{"inference_objects":[]}`. Response `200` memakai `data: null`.

## Master Model

### GET `/master/models`

Mengambil seluruh model tersimpan. Untuk kebutuhan frontend berdasarkan usecase, gunakan endpoint khusus di bawah.

### GET `/master/usecases/:usecase_code/models`

Mengambil semua model yang terhubung ke kombinasi inference-object yang dipakai usecase tersebut. Model yang sama dapat muncul pada beberapa usecase.

Contoh response:

```json
{
  "code": 200,
  "data": [
    {
      "id": 12,
      "file_url": "2e7c-model.onnx",
      "name": "People Detector",
      "model_name": "yolov8n",
      "version": "v1.0.0",
      "code": "YOLOV8N",
      "size": 12345678,
      "inference_object_id": 10,
      "inference_id": 1,
      "object_id": 2,
      "inference_name": "Detection",
      "inference_code": "detection",
      "object_name": "Person",
      "object_code": "person",
      "created_at": "2026-10-16T10:00:00Z"
    }
  ],
  "message": "Successfully get models by usecase"
}
```

### GET `/master/models/:model_id`

Mengambil detail model.

### GET `/master/models/:model_id/download`

Download file model berdasarkan `model_id`. Response berupa binary file dengan `Content-Disposition: attachment`.

### POST `/master/models`

Gunakan `multipart/form-data`:

| Field | Tipe | Wajib |
|---|---|---|
| `file` | file | Ya |
| `name` | string | Ya |
| `model_name` | string | Ya |
| `version` | string | Ya |
| `code` | string | Ya |
| `inference_object_ids` | string | Ya, comma-separated |

Response `201` mengembalikan object model.

### PUT `/master/models/:model_id`

Handler saat ini membaca body JSON:

```json
{
  "file_url": "model.onnx",
  "name": "People Detector",
  "model_name": "yolov8n",
  "version": "v1.0.0",
  "code": "YOLOV8N",
  "size": 12345678,
  "inference_object_ids": [10, 11]
}
```

Response `200` memakai `data: null`.

### DELETE `/master/models/:model_id`

Menghapus model. Response `200` memakai `data: null`.

### Endpoint lama

Endpoint berikut masih dipertahankan untuk kompatibilitas, tetapi bukan contract utama frontend:

- `GET /master/models/download/:usecase_code` — download model terbaru berdasarkan usecase code.
- `GET /master/models/download/inference-object/:inference_object_id` — download model terbaru berdasarkan inference-object.

## Capture proxy

### POST `/cv/capture`
### POST `/capture`

Kedua endpoint mengarah ke proxy yang sama dan meneruskan request ke CV service.

```json
{
  "session_id": 9,
  "usecase": {"id": 10, "code": "people_monitoring", "config": {"minimum_duration": 3}},
  "source": {"type": "stream", "fps": 30, "camera": {"camera_id": "CAM-001", "stream_address": "rtsp://192.168.1.10/live"}},
  "inference": {"type": "detection", "objects": ["person"]},
  "zones": [],
  "output": {"mode": "event", "interval": 1}
}
```

Proxy meneruskan status code, content type, dan body dari CV service. Error umum: `400` JSON invalid, `500` `BASE_URL_CV` belum dikonfigurasi, `502` CV service tidak dapat dihubungi.

## Source, camera, video, dan zone — belum exposed

Handler backend sudah tersedia, tetapi route belum dipasang di `internal/infrastructure/api/routes/master.go`. Request berikut akan menghasilkan `404` sampai route didaftarkan:

- `GET /master/sources`
- `POST /master/videos`
- `POST /master/cameras`
- `PUT /master/cameras/:camera_id`
- `GET /master/cameras/:camera_id/zones`
- `POST /master/cameras/:camera_id/zones`
- `PUT /master/cameras/:camera_id/zones/:zone_id`
- `DELETE /master/cameras/:camera_id/zones/:zone_id`

Kontrak request yang sudah disiapkan:

### Camera

```json
{
  "area_location_id": 1,
  "camera_code": "CAM-001",
  "camera_name": "Lobby Camera",
  "camera_status": true,
  "ip_address": "192.168.1.10",
  "stream_url": "rtsp://192.168.1.10/live",
  "auth_type": "basic",
  "username": "camera-user",
  "password": "camera-password",
  "installed_at": "2026-10-08T10:00:00Z",
  "usecase_config": [],
  "channels": "channel-1",
  "camera_type": "fixed",
  "capture_protocol": "rtsp",
  "capture_mode": "STREAM",
  "rtc_url": "https://example.local/camera/CAM-001",
  "device_path": null
}
```

Semua field camera opsional. `capture_mode` jika diisi hanya `STREAM` atau `SNAPSHOT`; `installed_at` RFC3339; password tidak pernah dikembalikan. Untuk USB camera gunakan `camera_type: "usb"` dan `device_path`.

### Zone

```json
{
  "zone_name": "Lobby",
  "zone_type": "polygon",
  "zone_coordinate": {"points": [[0.1, 0.1], [0.9, 0.1], [0.9, 0.9], [0.1, 0.9]]},
  "config": {"max_person_count": 6},
  "usecase_ids": [1, 2]
}
```

`zone_name`, `zone_type`, `zone_coordinate`, dan setiap ID pada `usecase_ids` wajib valid.

### Source list

```text
GET /master/sources?page=1&limit=10&search=camera&source_type=camera&start_time=2026-10-08T00:00:00Z&end_time=2026-10-09T00:00:00Z
```

`source_type` hanya `camera` atau `video`; waktu harus RFC3339; `start_time` inklusif dan `end_time` eksklusif.

## HTTP status

| Status | Arti |
|---:|---|
| 200 | Request berhasil |
| 201 | Resource berhasil dibuat |
| 400 | JSON, parameter, atau validasi salah |
| 401 | Token tidak ada atau session invalid |
| 404 | Resource/route tidak ditemukan |
| 500 | Error database/server |
| 502 | Proxy gagal menghubungi CV service |

## Local development

```bash
make sqlite.setup
make run
```

Health check:

```bash
curl http://127.0.0.1:8000/
```

Swagger tersedia saat `SWAGGER_ENABLED=true`, tetapi gunakan dokumen ini sebagai handoff yang mengikuti route runtime.
