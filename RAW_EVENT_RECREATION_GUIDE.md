# Raw Event Pipeline Recreation Guide

Blueprint untuk recreate alur `raw_event` camera service. Referensi utama: `raw_event_processing_optimization_brainstorm.md`.

## Source of truth implementasi

```text
transport/http/presenter/raw_event/create_raw_event.go
transport/http/handlers/raw_event/create_raw_event_handler.go
usecase/raw_event/create_raw_event_service.go
usecase/raw_event/service.go
usecase/raw_event/process_raw_event_service.go
repository/raw_event/create_raw_event_query.go
repository/raw_event/processing_query.go
repository/raw_event/repository.go
database/migrations/20260818000100_create_table_raw_event.up.sql
```

## 1. Flow ingress

```text
CV/client
  -> POST /api/v1/detection-events
  -> validate request
  -> resolve session_model
  -> insert raw_event (pending)
  -> resolve master_usecase.usecase_code
  -> enqueue raw_event_id
  -> HTTP 202 Accepted

worker
  -> claim pending/failed menjadi processing
  -> dispatch berdasarkan usecase_code
  -> insert detection_event/event domain
  -> processed atau failed
```

Ingress hanya melakukan validasi ringan dan insert durable. Processor berat berjalan asynchronous.

## 2. Request dan session

```go
type CreateRawEventRequest struct {
    Timestamp        string          `json:"timestamp" validate:"required"`
    FrameTimestamp   string          `json:"frame_timestamp"`
    FrameTimestampNS *int64          `json:"frame_timestamp_ns"`
    CameraID         int64           `json:"camera_id" validate:"required"`
    UsecaseID        int64           `json:"usecase_id"`
    SessionID        int64           `json:"session_id" validate:"required"`
    EventPayload     json.RawMessage `json:"event_payload" validate:"required"`
    ImageSnapshotURL string          `json:"image_snapshot_url"`
}
```

`session_id` adalah sumber relasi runtime. Ambil nilai efektif dari:

```sql
SELECT id, camera_id, usecase_id
FROM session_model
WHERE id = $1 AND deleted_at IS NULL;
```

Simpan `camera_id` dan `usecase_id` dari session, bukan mempercayai request jika berbeda. Response sukses berbentuk `raw_event_id + status=accepted` dengan HTTP `202`. Session/timestamp invalid menjadi `400`; enqueue gagal menjadi `503`.

## 3. Relasi master camera, zone, usecase

```text
master_camera
  -> master_camera_zone
      -> usecase_zone -> master_usecase

master_camera + master_usecase
  -> session_model
      -> raw_event
```

Foreign key raw event:

```text
raw_event.camera_id  -> master_camera.id
raw_event.usecase_id -> master_usecase.id
raw_event.session_id -> session_model.id
```

`master_camera` menyimpan identitas dan metadata capture. `master_camera_zone` menyimpan zone camera. `usecase_zone` membatasi usecase yang boleh memakai zone. `master_usecase.usecase_code` menentukan processor. `session_model` mengikat camera, usecase, dan scope konfigurasi.

Konfigurasi zone diambil dari `config_cv` dengan `session_id + zone_id`:

```sql
SELECT key, value
FROM config_cv
WHERE session_id = $1
  AND config_type = 'zone'
  AND type_id = $2
  AND deleted_at IS NULL;
```

Project baru sebaiknya memvalidasi `zone_id` terhadap camera dan `usecase_zone`.

### Master inference dan master object

Model inference tidak langsung ditempel ke raw event. Relasinya:

```text
master_usecase
  -> usecase_inference_object
       -> master_inference
       -> master_object
              |
              v
       session_model.inference_object_id
```

Schema inti:

```text
master_inference:
  id, inference_name, inference_code, description, deleted_at

master_object:
  id, object_name, object_code, description, deleted_at

usecase_inference_object:
  id, usecase_id, inference_id, object_id, deleted_at
```

`usecase_inference_object` adalah pasangan inference-object yang diizinkan untuk sebuah usecase. `session_model.inference_object_id` memilih pasangan yang dipakai oleh satu runtime session. Saat membuat session group, service mengambil semua pasangan aktif untuk camera + usecase, lalu membuat child session per pasangan.

Karena itu satu camera/usecase dapat memiliki beberapa session child, misalnya pasangan inference/object berbeda. Raw event tetap menunjuk satu `session_id`, sehingga processor dapat menentukan inference dan object secara deterministik.

## 4. Bentuk tabel dan foreign key

Bagian ini adalah bentuk logical schema berdasarkan migration aktif. Kolom audit (`created_*`, `updated_*`, `deleted_*`) dipertahankan karena service memakai soft delete. Project baru boleh menyesuaikan nama kolom, tetapi relasi dan invariant berikut harus tetap ada.

### 4.1 Tabel master camera dan usecase

```sql
CREATE TABLE master_camera (
  id BIGSERIAL PRIMARY KEY,
  area_location_id BIGINT,
  camera_code VARCHAR(50) NOT NULL,
  camera_name VARCHAR(255) NOT NULL,
  camera_status BOOLEAN NOT NULL DEFAULT TRUE,
  ip_address VARCHAR(45) NOT NULL,
  stream_url VARCHAR(255) NOT NULL,
  auth_type VARCHAR(50) NOT NULL DEFAULT 'basic',
  username VARCHAR(255),
  password VARCHAR(255),
  installed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

CREATE TABLE master_usecase (
  id BIGSERIAL PRIMARY KEY,
  usecase_name VARCHAR(255) NOT NULL,
  usecase_code VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

-- Migration berikutnya membuat usecase_code unik hanya untuk data aktif.
CREATE UNIQUE INDEX unique_usecase_code
  ON master_usecase (usecase_code) WHERE deleted_at IS NULL;
```

`master_camera` menyimpan identitas dan koneksi kamera. `master_usecase.usecase_code` adalah key dispatch worker; jangan dispatch berdasarkan nama display.

### 4.2 Zone camera dan izin zone untuk usecase

```sql
CREATE TABLE master_camera_zone (
  id BIGSERIAL PRIMARY KEY,
  camera_id BIGINT REFERENCES master_camera(id) ON DELETE CASCADE,
  zone_name VARCHAR(255) NOT NULL,
  zone_type VARCHAR(50) NOT NULL,
  zone_coordinate JSONB NOT NULL,
  config JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

CREATE TABLE usecase_zone (
  id BIGSERIAL PRIMARY KEY,
  usecase_id BIGINT REFERENCES master_usecase(id) ON DELETE CASCADE,
  master_camera_zone_id BIGINT
    REFERENCES master_camera_zone(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);
```

Validasi zone harus memastikan dua hal: `master_camera_zone.camera_id` sama dengan camera session, dan pasangan `(usecase_id, master_camera_zone_id)` aktif di `usecase_zone`.

### 4.3 Inference, object, dan pasangan yang diizinkan

```sql
CREATE TABLE master_inference (
  id BIGSERIAL PRIMARY KEY,
  inference_name VARCHAR(255) NOT NULL,
  inference_code VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

CREATE TABLE master_object (
  id BIGSERIAL PRIMARY KEY,
  object_name VARCHAR(255) NOT NULL,
  object_code VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

CREATE TABLE usecase_inference_object (
  id BIGSERIAL PRIMARY KEY,
  usecase_id BIGINT NOT NULL,
  inference_id BIGINT NOT NULL,
  object_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  UNIQUE (usecase_id, inference_id, object_id)
);
```

`usecase_inference_object` adalah tabel mapping, bukan tabel hasil inference. Karena migration saat ini hanya memberi index pada `inference_id` dan `object_id`, recreate yang lebih kuat sebaiknya menambahkan FK ke ketiga master tersebut dan index `usecase_id`.

### 4.4 Session group dan session runtime

```sql
CREATE TABLE session_model_group (
  id BIGSERIAL PRIMARY KEY,
  camera_id BIGINT NOT NULL,
  usecase_id BIGINT NOT NULL,
  output_mode VARCHAR(50) NOT NULL DEFAULT 'event',
  interval INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);

CREATE TABLE session_model (
  id BIGSERIAL PRIMARY KEY,
  usecase_id BIGINT NOT NULL,
  camera_id BIGINT NOT NULL,
  session_group_id BIGINT REFERENCES session_model_group(id),
  inference_object_id BIGINT
    REFERENCES usecase_inference_object(id) ON DELETE SET NULL,
  status BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100)
);
```

`session_model_group` menyimpan konfigurasi satu camera/usecase. `session_model` adalah runtime child; satu group dapat memiliki beberapa child, biasanya satu child per pasangan `inference_object_id`. Hanya session aktif dan belum dihapus yang boleh dipakai saat ingress.

### 4.5 Konfigurasi global dan konfigurasi runtime

```sql
CREATE TABLE usecase_config (
  id BIGSERIAL PRIMARY KEY,
  usecase_id BIGINT NOT NULL REFERENCES master_usecase(id),
  key VARCHAR(255) NOT NULL,
  value TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(100), updated_at TIMESTAMPTZ, updated_by VARCHAR(100),
  deleted_at TIMESTAMPTZ, deleted_by VARCHAR(100),
  UNIQUE (usecase_id, key)
);

CREATE TABLE config_cv (
  id BIGSERIAL PRIMARY KEY,
  session_id BIGINT NOT NULL,
  config_type VARCHAR(50) NOT NULL,
  type_id BIGINT NOT NULL,
  key VARCHAR(255) NOT NULL,
  value TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
  UNIQUE (session_id, config_type, type_id, key)
);
CREATE INDEX idx_config_cv_session_type
  ON config_cv (session_id, config_type, type_id);
```

Migration saat ini belum memasang FK `config_cv.session_id`; service tetap harus memvalidasi session sebelum insert. `type_id` ditafsirkan berdasarkan `config_type`: `usecase` → `master_usecase.id`, `zone` → `master_camera_zone.id`, dan `inference` → `usecase_inference_object.id` (fallback legacy: `master_inference.id`).

### 4.6 Raw event dan hasil turunan

```sql
CREATE TABLE raw_event (
  id BIGSERIAL PRIMARY KEY,
  camera_id BIGINT NOT NULL REFERENCES master_camera(id) ON DELETE RESTRICT,
  usecase_id BIGINT NOT NULL REFERENCES master_usecase(id) ON DELETE RESTRICT,
  session_id BIGINT NOT NULL REFERENCES session_model(id) ON DELETE RESTRICT,
  event_payload JSONB NOT NULL,
  event_timestamp TIMESTAMPTZ,
  frame_timestamp TIMESTAMPTZ,
  frame_timestamp_ns BIGINT,
  image_snapshot_url TEXT,
  processing_status VARCHAR(20) NOT NULL DEFAULT 'pending',
  processing_error TEXT,
  retry_count INTEGER NOT NULL DEFAULT 0 CHECK (retry_count >= 0),
  received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  processed_at TIMESTAMPTZ,
  CHECK (processing_status IN ('pending','processing','processed','failed'))
);

CREATE TABLE detection_event (
  id BIGSERIAL PRIMARY KEY,
  camera_id BIGINT REFERENCES master_camera(id) ON DELETE SET NULL,
  usecase_id BIGINT REFERENCES master_usecase(id) ON DELETE SET NULL,
  event_payload JSONB NOT NULL,
  asset_url VARCHAR(255) NOT NULL,
  event_time TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE detection_event_raw_event (
  id BIGSERIAL PRIMARY KEY,
  detection_event_id BIGINT NOT NULL
    REFERENCES detection_event(id) ON DELETE CASCADE,
  raw_event_id BIGINT NOT NULL
    REFERENCES raw_event(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (detection_event_id, raw_event_id)
);
```

Raw event adalah immutable input/audit log; `detection_event` adalah output processor. Satu raw event bisa menghasilkan beberapa detection event, dan relasinya dicatat di `detection_event_raw_event`.

### 4.7 Ringkasan relasi

```text
master_camera 1---N master_camera_zone
master_usecase 1---N usecase_zone
master_usecase 1---N usecase_inference_object
master_inference 1---N usecase_inference_object
master_object 1---N usecase_inference_object
session_model_group 1---N session_model
usecase_inference_object 1---N session_model
session_model 1---N config_cv
master_camera/master_usecase/session_model 1---N raw_event
raw_event N---N detection_event melalui detection_event_raw_event
```

Urutan pembuatan data yang aman: master camera/usecase → zone dan mapping usecase-zone → inference/object dan mapping pasangan → session group → child session → config CV → raw event.

## 5. Schema dan lifecycle

Kolom penting `raw_event`:

```text
id, camera_id, usecase_id, session_id
event_payload JSONB
event_timestamp, frame_timestamp, frame_timestamp_ns
image_snapshot_url
processing_status, processing_error, retry_count
received_at, processed_at
```

Lifecycle:

```text
pending -> processing -> processed
                    \-> failed -> processing
```

Insert:

```sql
INSERT INTO raw_event (
  camera_id, usecase_id, session_id, event_payload,
  event_timestamp, frame_timestamp, frame_timestamp_ns, image_snapshot_url
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING id;
```

Claim idempotent:

```sql
UPDATE raw_event
SET processing_status='processing', processing_error=NULL
WHERE id=$1 AND processing_status IN ('pending','failed')
RETURNING id;
```

Error menaikkan `retry_count` dan menyimpan `processing_error`. Sukses mengubah status menjadi `processed` dan mengisi `processed_at`.

Implementasi sekarang memakai buffered in-process channel berisi ID. Untuk multi-instance, gunakan DB-backed dispatcher/outbox dengan `FOR UPDATE SKIP LOCKED`.

## 6. Konfigurasi: usecase config dan config CV

Ada dua lapisan konfigurasi yang perlu dibedakan.

### `usecase_config`

`usecase_config` adalah konfigurasi dasar/global milik `master_usecase`:

```text
usecase_config:
  id, usecase_id, key, value, created_by, updated_by, deleted_at
```

Relasi:

```text
master_usecase 1 --- N usecase_config
```

Contoh key:

```text
max_person_count
min_count_change_time
min_status_hold_time
min_dwell_time
person_time_limit
face_recognition_threshold
```

Service usecase setting juga menyinkronkan konfigurasi ke legacy `master_camera.usecase_config` untuk kebutuhan compatibility. Project baru sebaiknya memilih satu source of truth; rekomendasinya PostgreSQL `usecase_config`, lalu cache effective config.

### `config_cv`

`config_cv` adalah konfigurasi runtime yang scoped ke session:

```text
config_cv:
  id, session_id, config_type, type_id, key, value,
  created_at, updated_at, deleted_at
```

Unique scope:

```text
(session_id, config_type, type_id, key)
```

Nilai `config_type` yang dipakai:

| `config_type` | Arti `type_id` | Contoh penggunaan |
|---|---|---|
| `usecase` | `usecase_id` | config umum session/usecase |
| `zone` | `master_camera_zone.id` | density, dwell, footfall, occupancy |
| `inference` | `usecase_inference_object.id` atau fallback `master_inference.id` | threshold/model behavior |

Saat session dibuat atau di-update, konfigurasi request ditulis ke `config_cv` untuk setiap child session. Konfigurasi inference dipetakan dari `inference_object_id` ke session yang sesuai.

### Effective config precedence

Untuk inference, implementasi sekarang memilih konfigurasi paling spesifik:

```text
1. config_cv.type_id = usecase_inference_object.id
2. fallback config_cv.type_id = master_inference.id
3. key yang tidak ditemukan dapat memakai row fallback umum
```

Untuk zone, processor membaca langsung:

```text
session_id + config_type='zone' + type_id=zone_id
```

Untuk usecase, gunakan:

```text
session_id + config_type='usecase' + type_id=usecase_id
```

Semua value disimpan sebagai `TEXT`, sehingga processor harus melakukan parsing typed (`int`, `float`, `bool`, enum) dengan default yang jelas.

## 7. Dispatch dari master usecase

```go
switch strings.ToUpper(usecaseCode) {
case "ZONE_DENSITY":
    return processZoneDensity(raw)
case "PEOPLE_DWELL_TIME":
    return processDwellTime(raw)
case "PEOPLE_FOOTFALL":
    return processFootfall(raw)
case "RECEPTIONIST_ANALYTICS", "ZONE_OCCUPANCY":
    return processReceptionistPeopleDetection(raw)
case "FACE_RECOGNITION":
    return processReceptionistFaceRecognition(raw)
default:
    return nil
}
```

Implementasi saat ini menganggap unknown usecase selesai tanpa event turunan. Project baru harus memilih apakah unknown usecase menjadi `failed` atau `ignored` dengan audit reason.

## 8. Processor `ZONE_DENSITY`

Payload umum:

```json
{"people_count":6,"zones":[{"zone_id":10,"people_count":6}]}
```

Flow: parse zone/count, baca config zone, tentukan `Crowded` jika count >= `max_person_count`, baca detection terakhir, skip perubahan terlalu cepat atau same-state, tunggu `min_status_hold_time`, lalu insert detection event per zone.

Default config saat tidak ada:

```text
max_person_count=6
min_count_change_time=5 detik
min_status_hold_time=15 detik
```

## 9. Processor `PEOPLE_DWELL_TIME`

Hanya action `exit` yang diproses. Cari raw `enter` terakhir dengan kombinasi camera/usecase/session/zone/person yang sama. Jika duration kosong, hitung dari enter dan exit. Default `min_dwell_time=15s` dan `person_time_limit=360s`; di bawah minimum di-skip, melewati limit menjadi `Overtime`, selainnya `Normal`. Insert detection event dan link ke raw enter serta raw exit.

Untuk traffic tinggi, gunakan tabel `person_dwell_sessions` agar tidak mencari raw history berulang.

## 10. Processor `PEOPLE_FOOTFALL`

Validasi action `cross` dan `zone_id`, cocokkan direction dengan config zone, pilih `crossing_timestamp`, lalu insert event turunan:

```json
{"zone_id":10,"total_people":1,"direction_type":"IN"}
```

Footfall adalah kandidat fast/batch lane.

## 11. Processor receptionist dan face recognition

`RECEPTIONIST_ANALYTICS` dan `ZONE_OCCUPANCY` memakai flow:

```text
raw people/zone payload
  -> active master_counter berdasarkan zone
  -> patient/staff counter atau face camera
  -> baca occupancy state terakhir
  -> stable: receptionist event
  -> unstable: receptionist_event_pending
```

Scheduler occupancy berjalan tiap 2 detik untuk stable candidates. `FACE_RECOGNITION` membaca employee/identity, confidence/similarity, threshold config, active counter, dan pending occupancy. External CV call harus dipisahkan ke worker/job tersendiri.

## 12. Detection event dan audit

`detection_event` adalah hasil turunan; raw event tetap disimpan.

```sql
INSERT INTO detection_event (camera_id,usecase_id,event_payload,asset_url,event_time)
VALUES ($1,$2,$3,$4,$5)
RETURNING id;

INSERT INTO detection_event_raw_event (detection_event_id,raw_event_id,created_at)
VALUES ($1,$2,NOW())
ON CONFLICT (detection_event_id,raw_event_id) DO NOTHING;
```

Satu raw event dapat menghasilkan beberapa detection event, misalnya satu per zone.

## 13. Snapshot cleanup

Snapshot yang tidak menghasilkan event turunan masuk cleanup queue dengan alasan seperti:

```text
zone_density_no_zones
zone_density_same_state
zone_density_min_status_hold_time
dwell_time_invalid_payload
dwell_time_enter_not_found
footfall_direction_mismatch
```

Scheduler menghapus object storage, menghapus queue jika sukses, dan retry jika gagal.

## 14. Rekomendasi recreate dan checklist

Fase functional parity:

1. buat master camera, camera zone, master usecase, usecase-zone, session model;
2. buat raw event JSONB dengan FK camera/usecase/session;
3. resolve session saat ingress;
4. insert raw event, enqueue ID, return 202;
5. claim idempotent dan dispatch berdasarkan usecase code;
6. implement density, dwell, footfall, receptionist, dan face processor;
7. buat detection event dan audit link.

Worker pool yang disarankan:

```text
ZONE_DENSITY      stateful pool + Redis transient state
ZONE_OCCUPANCY    occupancy pool + external CV pool
DWELL_TIME        ordered state pool
FOOTFALL          fast/batch pool
FACE_RECOGNITION  recognition pool
```

Ordering key:

```text
density: camera_id + usecase_id + zone_id
dwell:   camera_id + usecase_id + zone_id + person_id
```

Cache:

```text
session:model:{session_id}
config:{usecase_code}:session:{session_id}
```

Checklist minimum:

- [ ] master camera, camera zone, master usecase, usecase-zone, session model;
- [ ] raw event JSONB, FK, status, error, retry count;
- [ ] resolve camera/usecase dari session;
- [ ] HTTP 202 setelah insert/enqueue;
- [ ] worker claim idempotent dan stale-processing recovery;
- [ ] config zone berdasarkan `session_id + zone_id`;
- [ ] `usecase_config` sebagai konfigurasi dasar per master usecase;
- [ ] `config_cv` dengan scope `session_id + config_type + type_id + key`;
- [ ] master inference dan master object;
- [ ] mapping `usecase_inference_object`;
- [ ] `session_model.inference_object_id` untuk memilih konfigurasi inference;
- [ ] precedence config inference yang spesifik lalu fallback ke master inference;
- [ ] parsing value `TEXT` menjadi tipe config yang benar;
- [ ] validasi zone terhadap camera/usecase;
- [ ] detection event dan raw-event audit link;
- [ ] external CV worker terpisah;
- [ ] snapshot cleanup queue;
- [ ] ordering usecase stateful;
- [ ] metrics throughput, backlog, oldest pending age, processing duration;
- [ ] test ingress, lifecycle, skip rule, dan semua processor.

