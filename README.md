# Jualan Online - Sistem Manajemen Inventaris & Finansial Multi-Channel

Proyek ini adalah sistem manajemen inventaris dan finansial terpusat yang dirancang untuk membantu UMKM/Pemilik Usaha mengelola stok barang dan melacak profitabilitas di berbagai marketplace online (Shopee, Tokopedia, TikTok Shop, dll) dalam satu platform terpadu.

## 🚀 Fitur Utama

- **Sentralisasi Inventaris:** Sinkronisasi stok barang untuk semua channel penjualan.
- **Manajemen Pembelian (Procurement):** Pelacakan pembelian barang dari supplier/marketplace asal untuk menentukan COGS (*Landed Cost*).
- **Audit Pergerakan Stok:** Pencatatan otomatis setiap barang masuk (pembelian), keluar (penjualan), dan retur.
- **Laporan Finansial:** Perhitungan margin keuntungan dan profit/loss per transaksi secara real-time.
- **Dukungan Multi-Channel:** Manajemen toko untuk berbagai platform marketplace.

## 🏗️ Arsitektur Sistem

Sistem ini dibangun menggunakan arsitektur **Microservices** yang bersifat *event-driven*:

- **Backend:** Go (Golang) menggunakan framework Echo.
- **Frontend:** Next.js 16 (App Router), TailwindCSS, shadcn/ui.
- **Persistence:** PostgreSQL (Isolasi database per service).
- **Messaging:** Redis Streams untuk komunikasi antar-layanan asinkron yang andal.
- **API Gateway:** Traefik untuk routing dan SSL termination.
- **Deployment:** Docker Swarm dengan infrastruktur berbasis Docker Compose untuk pengembangan lokal.

## 🛠️ Tech Stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL 15
- **Caching & Message Broker:** Redis 7 (Streams)
- **API Gateway:** Traefik v3
- **Containerization:** Docker & Docker Compose

## 📂 Struktur Proyek

```text
.
├── services/
│   ├── common/              # Shared library internal (Logger, DB, Messaging, Config)
│   ├── auth-service/        # Service Autentikasi & User Management
│   └── (service lainnya)    # Inventory, Sales, Procurement, dll.
├── docker-compose.yml       # Orkestrasi layanan lokal
├── Makefile                 # Shortcut perintah pengembangan
├── PRD.md                   # Dokumen kebutuhan produk (Source of Truth)
├── GEMINI.md                # Panduan instruksional pengembangan
└── .env.example             # Template variabel lingkungan
```

## 🏁 Memulai Pengembangan (Local Setup)

### Prasyarat
- Docker & Docker Compose
- Go 1.22+
- Make (opsional, untuk menjalankan Makefile)

### Langkah-langkah
1.  **Clone Repositori:**
    ```bash
    git clone https://github.com/maztea/jualan-online.git
    cd jualan-online
    ```
2.  **Siapkan Env:**
    ```bash
    cp .env.example .env
    ```
3.  **Jalankan Infrastruktur:**
    ```bash
    make up
    ```
4.  **Seed Data Awal (User Admin/Staff):**
    ```bash
    make seed
    ```

## 🔄 Alur Kerja Git (Git Workflow)

Sesuai dengan panduan di `GEMINI.md`, proyek ini mewajibkan penggunaan **`git worktree`** untuk setiap implementasi fitur baru:
1.  Buat worktree baru untuk fitur: `git worktree add -b feature/nama-fitur .worktrees/nama-fitur`
2.  Implementasikan kode di direktori worktree tersebut.
3.  Komit dan gabungkan kembali ke `main`.

---
*Proyek ini masih dalam tahap pengembangan aktif.*
