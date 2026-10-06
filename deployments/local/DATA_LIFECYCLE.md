# Local Stack — Data Lifecycle & Recovery

Dokumen ini menjelaskan cara data dikelola di local stack (`deployments/local/`)
sesuai checklist **Fase 6 — Local Docker Compose dan Data Lifecycle**.

Ringkasan cepat:

| Topik | Lokasi |
|---|---|
| Compose file | `deployments/local/docker-compose.yml` |
| Env (dev) | `deployments/local/docker.env` (+ `.example`) |
| Logs | `deployments/local/logs/*.log` (di-scrape Promtail → Loki) |
| Backup DB | `./scripts/backup.sh` → `deployments/local/backups/` |
| Restore DB | `./scripts/restore.sh` |
| Smoke test | `./scripts/smoke-test.sh` |
| Migrations | `service/migrate/` (Goose), `just migrate` / `just migrate-down` |

---

## 1. Topologi & penyimpanan data

State yang bertahan (persist) didefinisikan sebagai named volumes di compose:

| Volume | Data |
|---|---|
| `postgres_identity_data`, `postgres_merchant_data`, `postgres_catalog_data`, `postgres_sales_data`, `postgres_experience_data`, `postgres_email_data` | PostgreSQL `/var/lib/postgresql/data` — satu volume per bounded context |
| `kafka_ecommerce_data` | Kafka `/bitnami/kafka` |
| `redis_ecommerce_data` | Redis `/data` |
| `clickhouse_ecommerce_data` | ClickHouse `/var/lib/clickhouse` |
| `loki_ecommerce_data` | Log Loki |
| `grafana_ecommerce-storage` | Dashboard Grafana |

Volume ini **tidak dihapus** oleh `docker compose down` (hanya `down -v` yang
menghapusnya). Selama diagnosa normal gunakan `docker compose down` (tanpa `-v`).

```bash
# Hentikan stack, data tetap aman
docker compose -f deployments/local/docker-compose.yml down

# Mulai ulang stack dari data yang sama
docker compose -f deployments/local/docker-compose.yml up -d
```

> ⚠️ **Hati-hati:** `docker compose down -v` menghapus semua volume. Hanya
> gunakan bila memang ingin reset total (lihat §6).

---

## 2. Backup PostgreSQL

Database-nya enam instance terpisah — satu per bounded context (`identity`,
`merchant`, `catalog`, `sales`, `experience`, `email`) — jadi satu kali backup
menghasilkan satu snapshot berisi enam dump.

```bash
# Manual: simpan ke deployments/local/backups/ecommerce_<ts>/
./deployments/local/scripts/backup.sh

# Dengan retention 14 hari
./deployments/local/scripts/backup.sh 14
```

Yang dilakukan script:
1. Memastikan keenam container `postgres_<ctx>` sehat.
2. `pg_dump --no-owner --no-privileges` per context → gzip →
   `backups/ecommerce_<ts>/ec_<ctx>.sql.gz`, plus file `MANIFEST`.
3. Menghapus snapshot lebih tua dari retention (default 7 hari).

Hasil dump **portable** (tanpa owner/privilege) sehingga bisa di-restore ke
environment sekali pakai (disposable).

---

## 3. Restore PostgreSQL

```bash
# Restore dengan konfirmasi
./deployments/local/scripts/restore.sh backups/ecommerce_20260101_000000

# Restore tanpa prompt (untuk CI/script)
./deployments/local/scripts/restore.sh --yes backups/ecommerce_20260101_000000
```

Yang dilakukan script:
1. Memvalidasi direktori snapshot, tiap dump (non-kosong + gzip valid), dan
   container `postgres_<ctx>` untuk setiap context di dalamnya.
2. Men-terminate koneksi aktif, lalu `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`
   di tiap database context.
3. Meng-restore tiap dump `ec_<ctx>.sql.gz`.

> ⚠️ Restore bersifat **destruktif** — schema `public` di setiap database
> context di-drop dan dibuat ulang. Selalu pastikan backup terbaru tersedia
> sebelum restore.

Schema di-drop (bukan database-nya) supaya target PgBouncer tetap valid dan
koneksi service yang sedang berjalan bisa reconnect tanpa restart.

**Restore ke environment disposable:** karena dump memakai `--no-owner
--no-privileges`, cukup jalankan restore terhadap postgres baru di env lain —
tidak perlu user/role khusus.

---

## 4. Migration & rollback (Goose)

Migration ada di `service/<svc>/database/migration/` (satu direktori per service;
service yang berbagi bounded context berbagi satu database) dan dijalankan dengan
Goose.

```bash
# Migrate ke versi terbaru
just migrate                        # go run service/migrate/cmd/main.go up

# Status migration
go run service/migrate/cmd/main.go status

# Rollback satu versi
just migrate-down                   # go run service/migrate/cmd/main.go down

# Rollback ke versi tertentu (mis. 00042)
go run service/migrate/cmd/main.go down-to 42

# Rollback total (kosongkan semua migration)
go run service/migrate/cmd/main.go reset
```

**Prosedur jika migration gagal:**

1. Baca output `go run service/migrate/cmd/main.go up` — tiap context di-log
   dengan prefix `[identity]`, `[merchant]`, dst.
2. Identifikasi versi yang gagal: `go run service/migrate/cmd/main.go status`.
3. Rollback ke versi terakhir yang sukses: `go run service/migrate/cmd/main.go down-to <n>`.
4. Perbaiki file migration, lalu `just migrate` lagi.
5. Pastikan service yang bergantung pada kolom baru di-restart setelah migrate.

> Goose mencatat versi terakhir yang berhasil di tabel `goose_db_version` —
> `down`/`down-to` aman karena hanya membalik migrasi terdaftar. Tabel itu
> dimiliki per bounded context, dan service yang berbagi satu context
> (mis. `auth`, `user`, `role` di identity) berbagi satu urutan versi.

---

## 5. Verifikasi stack (cold start & restart)

Gunakan smoke test setelah **cold start** (volume baru) dan setelah **restart**
(data lama) untuk memastikan state bertahan dan endpoint sehat.

```bash
# 1. Build image + start stack
just build-up

# 2. Tunggu semua service healthy, lalu smoke test
./deployments/local/scripts/smoke-test.sh

# 3. Restart (data harus bertahan — tanpa -v)
docker compose -f deployments/local/docker-compose.yml restart apigateway
./deployments/local/scripts/smoke-test.sh
```

Smoke test memverifikasi (bukan sekadar container "running"):
- `GET /health` → 200
- `GET /ready` → 200 (readiness cek Redis)
- `POST /api/auth/register` + `POST /api/auth/login` → token
- `GET /api/auth/me` (JWT) → 200
- `GET /api/product-query` (JWT) → 200

---

## 6. Reset total (opsional)

Gunakan hanya bila ingin membuang semua state:

```bash
# Hentikan + hapus volume
docker compose -f deployments/local/docker-compose.yml down -v

# Jalankan migrasi + seeder dari awal
just migrate
just seeder

# Verifikasi
./deployments/local/scripts/smoke-test.sh
```

Backup terlebih dahulu (§2) bila masih ada data yang dibutuhkan.

---

## 7. Retention & diagnosa non-destruktif

| Aset | Retention | Catatan |
|---|---|---|
| DB backup (`backups/`) | 7 hari (configurable) | Script backup mem-prune otomatis |
| Log service (`logs/`) | sesuai kebijakan lokal | Jangan dihapus saat diagnosa; di-scrape Promtail |
| Loki data (volume) | sesuai konfigurasi `loki-config.yaml` | Retention default local |
| Kafka/Redis volume | selama stack ada | Dihapus hanya via `down -v` |

Aturan diagnosa:

- Jangan hapus volume saat menyelidiki masalah — gunakan `logs` + `/health` +
  `/ready` + Prometheus/Grafana.
- Untuk "service tidak sehat" cek **readiness** (`/ready`), bukan hanya status
  container.
- Untuk masalah data korup: backup dulu, baru reset service tertentu (bukan
  seluruh volume).

---

## 8. Secret dev vs produksi

- `deployments/local/docker.env` — rahasia development (sudah ada + `.example`).
- Jangan pernah memakai nilai dev di produksi.
- `SECRET_KEY` dev (dipakai JWT gateway) berbeda dari secret produksi.
- SMTP credential dev (Mailtrap dsb.) terpisah dari SMTP produksi.
