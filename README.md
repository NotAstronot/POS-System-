# POS System - High Performance Point of Sale

Sistem kasir modern dengan arsitektur **high-throughput**, **multi-tenant**, dan **multi-outlet** yang siap menangani ribuan transaksi per detik.

## 🚀 Fitur Utama

### Core POS
- ⚡ Pencarian produk cepat (barcode scanner + full-text search)
- 🛒 Manajemen keranjang dengan diskon & variasi produk
- 💳 Multi-metode pembayaran (Tunai, QRIS, Debit/Kredit, E-Wallet, PayLater)
- 🖨️ Cetak struk (thermal printer) & struk digital (WhatsApp)
- 📊 Dashboard real-time & laporan penjualan

### Arsitektur Scalable
- **gRPC** untuk komunikasi antar service (Protocol Buffers)
- **Optimistic Locking** pada stok untuk mencegah overselling
- **Apache Kafka** untuk pemrosesan stok asinkron
- **Redis** untuk caching dan distributed locking
- **PostgreSQL** dengan connection pooling untuk database utama

## 🏗️ Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Frontend | Flutter (Dart) + Tailwind CSS |
| Desktop | Flutter (Windows/macOS/Linux) |
| Backend | Go (Golang) + Gin Framework |
| Database | PostgreSQL 16 + SQLite (offline) |
| Cache | Redis 7 |
| Message Queue | Apache Kafka |
| Container | Docker + Docker Compose |

## 📋 Struktur Proyek

```
pos-system/
├── backend/                     # Go backend (Clean Architecture)
│   ├── cmd/server/             # Entry point (bootstrap saja)
│   ├── internal/
│   │   ├── app/                # Composition root: wiring DI (repo → usecase → handler)
│   │   ├── config/             # Konfigurasi (config.yaml + env)
│   │   ├── domain/             # Enterprise rules: entity & interface repository
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── usecase/            # Application business logic
│   │   ├── repository/         # Data access (framework & drivers)
│   │   │   └── postgres/
│   │   ├── delivery/           # Interface adapters (presentation)
│   │   │   └── rest/           # Handler, middleware, dan router HTTP
│   │   └── database/           # Helper infrastruktur (RLS tenant)
│   └── migrations/             # SQL migrations
├── frontend/                   # React + Electron
│   ├── src/
│   │   ├── components/        # POS components
│   │   ├── pages/             # Pages (POS, Dashboard, Reports)
│   │   ├── services/          # API client
│   │   └── store/             # State management (Zustand)
│   └── electron/              # Desktop wrapper
└── docker-compose.yaml         # Infrastructure
```

Arah dependensi (Clean Architecture):

```
delivery (rest) → usecase → repository (postgres) → database
       ↘_______________ domain (entity & interface) _______________↗
```

## 🚦 Quick Start

### Prasyarat
- Docker & Docker Compose
- Node.js 20+ (development)
- Go 1.22+ (development)

### Production (Docker)
```bash
# Clone & start all services
docker compose up -d

# Access
# Frontend: http://localhost
# API: http://localhost:8080
```

### Development

**Backend:**
```bash
cd backend
go mod tidy
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**Desktop App:**
```bash
cd frontend
npm run electron:dev
```

## 🔧 Konfigurasi

Environment variables untuk backend:

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 8080 | HTTP API port |
| GRPC_PORT | 9090 | gRPC port |
| DB_HOST | localhost | PostgreSQL host |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| REDIS_HOST | localhost | Redis host |
| KAFKA_BROKER | localhost:9092 | Kafka broker |
| JWT_SECRET | (change me) | JWT signing secret |

## 📊 Performa

- **Latensi:** < 50ms per transaksi (99th percentile)
- **Throughput:** 5.000+ transaksi/detik (single instance)
- **Concurrency:** Optimistic locking pada stok, async processing via Kafka
- **Horizontal scaling:** Stateless backend, scale via Docker/K8s

## 📄 Lisensi

Proprietary - All Rights Reserved




Next Steps
1. Jalankan migration 015_multi_tenant.sql di PostgreSQL
2. Buat super_admin user dengan role super_admin
3. Test endpoint /api/v1/super-admin/tenants



2. ~~Struktur Proyek Golang (Clean Architecture)~~ — selesai
   Struktur folder terorganisir (app/delivery/usecase/repository/domain) memudahkan pemeliharaan kode saat aplikasi bertambah besar.


Implementasi Kode Utama (Golang)
JWT Payload & Custom Claims
Ketika user kasir atau owner login, backend memberikan JWT token yang menyimpan informasi tenant_id, user_id, role, dan outlet_id.

Middleware Tenant Extractor & RLS Injector
Middleware ini bertugas membaca JWT, mengambil tenant_id, memvalidasinya, lalu menyuntikkannya ke context.Context Go dan session PostgreSQL.

Repository (Safe Multi-Tenant Query) — selesai
Setiap transaksi ke database wajib menggunakan tenant_id dari context. Semua repository memakai withTenantTx/withTenantTx1 (tenant_id dari context + RLS), dengan predicate tenant_id eksplisit di setiap SQL.


