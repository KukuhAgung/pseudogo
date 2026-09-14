# PseudoGo

Web app untuk mengubah pseudocode (sesuai standar penulisan Algoritma Pemrograman kampus) menjadi kode Go yang valid — dan bisa langsung dijalankan, bukan cuma dilihat kodenya.

## Arsitektur

```
Browser (Angular)
      |
      v
Backend (Go) ── /convert  → parse + generate kode Go
      |         /run      → parse + generate + eksekusi
      v
Piston (Docker, terpisah dari repo ini) → sandbox eksekusi kode
```

Repo ini berisi dua bagian:

```
pseudogo/
├── backend/    Go — lexer, parser, code generator, REST API
└── frontend/   Angular — code editor, tampilan convert & run
```

**Piston tidak ada di dalam repo ini.** Ini layanan eksekusi kode pihak ketiga (open-source, self-hosted) yang dijalankan terpisah lewat Docker — mirip seperti database yang dijalankan di luar source code aplikasi. Backend terhubung ke Piston lewat HTTP di alamat yang bisa diatur (lihat bagian [Environment Variables](#environment-variables)).

## Prasyarat

Install semua ini terlebih dahulu:

| Tool | Versi minimum | Untuk |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.22+ | Backend |
| [Node.js](https://nodejs.org/) + npm | LTS terbaru | Frontend (Angular) |
| [Docker Desktop](https://www.docker.com/products/docker-desktop/) | terbaru | Menjalankan Piston |
| [Git](https://git-scm.com/) | — | Clone repo |

> Windows: pastikan virtualization sudah aktif di BIOS/UEFI (Docker Desktop akan memberi tahu kalau belum, lihat [Task Manager > Performance > CPU > Virtualization](#)).

## 1. Clone repo

```bash
git clone https://github.com/kukuhagung/pseudogo.git
cd pseudogo
```

## 2. Jalankan Piston (sandbox eksekusi kode)

Piston perlu dijalankan **sekali** dan akan tetap hidup di background (flag `-d`), terlepas dari backend/frontend dinyalakan atau tidak.

Buat folder khusus di luar repo ini (Piston bukan bagian dari source code proyek):

```bash
mkdir ~/piston-data && cd ~/piston-data
```

**Linux / macOS:**
```bash
docker run -d \
  --privileged \
  --restart unless-stopped \
  -v "$PWD:/piston" \
  -p 2000:2000 \
  --name piston_api \
  -e PISTON_RUN_TIMEOUT=30000 \
  -e PISTON_RUN_CPU_TIME=20000 \
  ghcr.io/engineer-man/piston
```

**Windows (PowerShell):**
```powershell
docker run -d `
  --privileged `
  --restart unless-stopped `
  -v "${PWD}:/piston" `
  -p 2000:2000 `
  --name piston_api `
  -e PISTON_RUN_TIMEOUT=30000 `
  -e PISTON_RUN_CPU_TIME=20000 `
  ghcr.io/engineer-man/piston
```

Piston baru nyala tanpa bahasa pemrograman apa pun terpasang. Install runtime Go-nya lewat CLI khusus (perlu Node.js — ini clone repo *terpisah*, tidak ada hubungannya dengan repo PseudoGo ini, cuma numpang lewat untuk instalasi):

```bash
git clone https://github.com/engineer-man/piston
cd piston/cli
npm install
node index.js ppman install go
```

Verifikasi:
```bash
curl http://localhost:2000/api/v2/runtimes
```
Harus muncul entry `"language": "go"` di responsnya.

## 3. Jalankan backend

```bash
cd backend
go mod download
go test ./...   # memastikan semuanya benar sebelum lanjut
go run ./cmd/server
```

Backend jalan di `http://localhost:8080`.

## 4. Jalankan frontend

Di terminal terpisah:
```bash
cd frontend
npm install
ng serve
```

Buka `http://localhost:4200` di browser.

## Environment Variables

Backend membaca dua environment variable opsional (kalau tidak diset, pakai default untuk development lokal):

| Variable | Default | Kegunaan |
|---|---|---|
| `PISTON_URL` | `http://localhost:2000/api/v2/execute` | Alamat API Piston |
| `ALLOWED_ORIGIN` | `http://localhost:4200` | Origin yang diizinkan CORS (alamat frontend) |

Set lewat environment sebelum menjalankan backend, contoh:
```bash
PISTON_URL="http://localhost:2000/api/v2/execute" ALLOWED_ORIGIN="http://localhost:4200" go run ./cmd/server
```

## Ringkasan urutan menjalankan

Setiap kali mau development, jalankan tiga hal ini di tiga terminal terpisah (urutan penting — Piston & backend harus nyala dulu sebelum frontend dipakai):

1. Piston sudah jalan di background (cukup sekali, tidak perlu diulang tiap sesi kecuali komputer restart)
2. `cd backend && go run ./cmd/server`
3. `cd frontend && ng serve`

## Kontribusi

Lihat [contributing.md](./contributing.md) dan [roadmap.md](./roadmap.md) untuk rencana pengembangan lanjutan (termasuk fitur OCR pseudocode tulisan tangan yang masih dalam tahap perencanaan).

## Lisensi

[MIT License](./LICENSE)