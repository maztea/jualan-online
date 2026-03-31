# Product Requirement Document (PRD): Sistem Manajemen Inventaris & Finansial Multi-Channel

## 1. EXECUTIVE SUMMARY

### Problem Statement & Goals
Pemilik usaha penjualan online menghadapi tantangan dalam sinkronisasi stok di berbagai marketplace (Shopee, TikTok Shop, Tokopedia, Blibli, Lazada) dan kesulitan melacak profitabilitas secara akurat. Tanpa pencatatan modal (COGS) yang terintegrasi, sulit untuk menentukan apakah operasional menghasilkan profit bersih.

**Goals:**
- Sentralisasi pengelolaan stok untuk semua channel penjualan.
- Otomasi perhitungan modal (COGS) dan margin keuntungan per transaksi.
- Laporan finansial dan inventaris yang komprehensif dan *real-time*.
- Audit trail untuk setiap pergerakan barang (pembelian, penjualan, retur).

### Target User Persona
- **Owner/Admin:** Mengelola toko, memantau laporan profit/loss, dan mengelola akses user.
- **Warehouse Staff:** Mencatat penerimaan barang, mengupdate stok fisik, dan memproses retur.
- **Purchasing Officer:** Mencatat pembelian barang dari supplier/marketplace lain untuk stok.

---

## 2. MICROSERVICES ECOSYSTEM & ARCHITECTURE

### Service Decomposition
Sistem akan dibangun dengan arsitektur microservices untuk skalabilitas dan independensi deployment:

1.  **Auth Service:** Mengelola user, role (RBAC), dan otentikasi.
2.  **Store Service:** Registrasi dan konfigurasi toko/marketplace (Shopee, Tokopedia, dll).
3.  **Inventory Service:** Core engine untuk tracking stok, katalog produk, dan log pergerakan barang.
4.  **Procurement Service:** Mengelola alur pembelian barang (Purchase Order) dan penerimaan barang (Goods Receipt) untuk menentukan COGS.
5.  **Sales Service:** Mencatat penjualan barang dan integrasi dengan pengurangan stok.
6.  **Finance & Report Service:** Mengagregasi data dari Procurement dan Sales untuk laporan profit/loss dan valuasi stok.

### Communication Pattern
- **Synchronous (REST API):** Digunakan untuk operasi CRUD sederhana dan interaksi user ke sistem via API Gateway.
- **Asynchronous (Event-Driven via Redis Streams):** Digunakan untuk sinkronisasi antar-layanan dengan jaminan persistensi pesan melalui *Consumer Groups*. Contoh: Saat `Sales Service` mencatat penjualan, event `order.created` dikirim ke stream untuk dikonsumsi oleh `Inventory Service` guna memotong stok.
- **Caching (Redis):** Digunakan di `Inventory Service` untuk akses cepat data stok terkini.

### Tech Stack
- **Frontend:** Next.js 16 (App Router), TailwindCSS, shadcn/ui.
- **Backend:** Go (Golang) dengan framework Gin.
- **Persistence:** PostgreSQL (1 per service untuk data isolation).
- **Caching & Message Broker:** Redis (Streams for Messaging).
- **API Gateway:** Traefik.

---

## 3. FUNCTIONAL REQUIREMENTS (User Stories)

| ID | User Story | Acceptance Criteria |
| :--- | :--- | :--- |
| **US-01** | As an Admin, I want to manage online stores so that I can categorize my sales channels. | 1. CRUD Store (Name, Platform Type, Status). <br> 2. Support platforms: Shopee, TikTok Shop, Tokopedia, Blibli, Lazada, Alfagift. |
| **US-02** | As a Purchasing Officer, I want to record purchases from marketplaces so that I can track my capital. | 1. Form input: Store Origin, Item, Qty, Unit Price, Shipping Cost. <br> 2. Calculate Total Cost per Item (Landed Cost). |
| **US-03** | As a Warehouse Staff, I want to record Goods Receipt so that stock is automatically updated. | 1. Match receipt with Purchase Order. <br> 2. Auto-increment stock in `Inventory Service`. <br> 3. Update Weighted Average Cost (WAC) for COGS calculation. |
| **US-04** | As an Admin, I want to see a Profit/Loss report so that I know my business health. | 1. Display Revenue - COGS - Expenses. <br> 2. Filter by Date Range and Store. |
| **US-05** | As a Staff, I want to record returns so that stock is adjusted correctly. | 1. Input return reason (Damaged/Cancel). <br> 2. Option to return to "Good Stock" or "Rejected Stock". |
| **US-06** | As a Staff, I want to record a sale/order from a specific store so that I can track revenue. | 1. Input Store ID, Marketplace Order ID, Item, Qty, Selling Price, and Marketplace Fees. <br> 2. Trigger `order.created` event to update inventory. |
| **US-07** | As a Staff, I want to view sales transaction history so that I can monitor daily store performance. | 1. List view of all orders. <br> 2. Filter by Store, Date, and Marketplace Order ID. |
| **US-08** | As an Admin, I want to see the gross margin per order so that I can evaluate product profitability. | 1. Calculate Margin = (Selling Price - COGS - Marketplace Fees - Shipping Discount). |

---

## 4. TECHNICAL SPECIFICATIONS & DATA CONTRACT

### Database Schema (PostgreSQL Samples)

**Auth Service:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    password_hash TEXT,
    role VARCHAR(20), -- ADMIN, STAFF
    created_at TIMESTAMP
);
```

**Store Service:**
```sql
CREATE TABLE stores (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    platform_type VARCHAR(50), -- SHOPEE, TOKOPEDIA, etc
    is_active BOOLEAN DEFAULT true
);
```

**Inventory Service:**
```sql
CREATE TABLE products (
    id UUID PRIMARY KEY,
    sku VARCHAR(50) UNIQUE,
    name VARCHAR(255),
    current_stock INT DEFAULT 0,
    average_cost DECIMAL(15,2) -- Used for COGS
);

CREATE TABLE stock_logs (
    id UUID PRIMARY KEY,
    product_id UUID REFERENCES products(id),
    change_qty INT,
    reference_type VARCHAR(50), -- PURCHASE, SALE, RETURN, ADJUSTMENT
    reference_id UUID,
    created_at TIMESTAMP
);
```

**Procurement Service:**
```sql
CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY,
    source_store_id UUID,
    status VARCHAR(20), -- DRAFT, ORDERED, RECEIVED, CANCELLED
    total_amount DECIMAL(15,2),
    shipping_cost DECIMAL(15,2),
    created_at TIMESTAMP
);

CREATE TABLE purchase_order_items (
    id UUID PRIMARY KEY,
    purchase_order_id UUID REFERENCES purchase_orders(id),
    product_id UUID,
    qty_ordered INT,
    qty_received INT DEFAULT 0,
    unit_price DECIMAL(15,2),
    landed_cost_per_unit DECIMAL(15,2)
);

CREATE TABLE goods_receipts (
    id UUID PRIMARY KEY,
    purchase_order_id UUID REFERENCES purchase_orders(id),
    received_by UUID,
    received_at TIMESTAMP,
    notes TEXT
);
```

**Sales Service:**
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    store_id UUID,
    marketplace_order_id VARCHAR(100),
    total_selling_price DECIMAL(15,2),
    marketplace_fees DECIMAL(15,2),
    status VARCHAR(20), -- COMPLETED, CANCELLED
    created_at TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    product_id UUID,
    qty INT,
    unit_price DECIMAL(15,2),
    cogs_at_time_of_sale DECIMAL(15,2)
);
```

### API Contracts (JSON)

**1. Auth Service: Login**
`POST /v1/auth/login`
```json
{
  "username": "admin",
  "password": "securepassword"
}
```

**2. Store Service: Create Store**
`POST /v1/stores`
```json
{
  "name": "My Shopee Store",
  "platform_type": "SHOPEE"
}
```

**3. Procurement Service: Receive Goods**
`POST /v1/procurement/orders/{id}/receive`
```json
{
  "received_items": [
    {
      "product_id": "uuid-string",
      "qty_received": 20
    }
  ],
  "received_by": "user-uuid-string"
}
```

**4. Inventory Service: Create Product**
`POST /v1/products`
```json
{
  "sku": "PROD-001",
  "name": "Barang Contoh A",
  "initial_stock": 0
}
```

**5. Inventory Service: Get Stock & COGS**
`GET /v1/products/{id}`
Response:
```json
{
  "id": "uuid-string",
  "sku": "PROD-001",
  "name": "Barang Contoh A",
  "current_stock": 50,
  "average_cost": 45000
}
```

**6. Inventory Service: Stock Adjustment (Return/Manual)**
`POST /v1/inventory/adjust`
```json
{
  "product_id": "uuid-string",
  "adjustment_qty": 5,
  "reason": "RETURN_GOOD_CONDITION",
  "reference_type": "SALES_RETURN",
  "reference_id": "order-uuid-string"
}
```

**7. Finance & Report Service: Profit/Loss Summary**
`GET /v1/reports/profit-loss?start_date=2023-01-01&end_date=2023-01-31`
Response:
```json
{
  "total_revenue": 15000000,
  "total_cogs": 10000000,
  "total_marketplace_fees": 500000,
  "gross_profit": 4500000
}
```

---

## 5. NON-FUNCTIONAL & DEVOPS REQUIREMENTS

### Scalability
- **Horizontal Scaling:** Deployment menggunakan Docker Swarm. Setiap service didefinisikan sebagai `service` di Swarm dengan `replicas: 3`.
- **Message Scalability:** Redis Streams dengan *Consumer Groups* untuk mendistribusikan beban kerja antar-replika service.

### Security
- **JWT Authentication:** Token dikelola oleh `Auth Service`.
- **RBAC:** Klaim role dalam JWT (`admin`, `staff`) divalidate di level API Gateway (Traefik Forward Auth) atau middleware tiap service.
- **HTTPS:** Enkripsi TLS untuk semua traffic eksternal dikelola oleh Traefik.

### Observability
- **Logging:** Structured logging menggunakan Zap (Go) dikirim ke ELK Stack atau Loki.
- **Monitoring:** Prometheus eksportir di tiap service, divisualisasikan dengan Grafana.
- **Tracing:** Jaeger untuk melacak latency antar-service saat proses sinkronisasi stok.

---

## 6. INFRASTRUCTURE & DEPLOYMENT

### Dockerization
Multi-stage build untuk efisiensi image size:
```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### CI/CD Pipeline
1.  **Linting & Test:** Menjalankan `golangci-lint` dan `go test ./...` pada setiap Push/PR.
2.  **Build Image:** Build Docker image dan push ke private registry.
3.  **Deploy to Swarm:** Menggunakan `docker stack deploy` melalui GitLab CI atau GitHub Actions untuk mengupdate service di environment production tanpa downtime (*Rolling Update*).
