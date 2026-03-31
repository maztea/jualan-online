# Project Context: Jualan Online Plan

This directory serves as the **Planning and Architectural Design** phase for the "Sistem Manajemen Inventaris & Finansial Multi-Channel". It contains the foundational requirements and microservices architecture specifications.

## 📁 Directory Overview
This project is currently in the **discovery and design stage**. It holds research notes, business requirements (PRD), and architectural blueprints for a system designed to centralize inventory and financial tracking across multiple online marketplaces (Shopee, Tokopedia, TikTok Shop, etc.).

## 🔑 Key Files
- **`PRD.md`**: The core Product Requirement Document. It defines the problem statement, user personas, functional requirements (user stories), and technical specifications.
- **`GEMINI.md`**: **MANDATORY READING**. Berisi pola implementasi dan aturan baku yang harus diikuti untuk setiap pengembangan fitur atau service baru.

## 🏗️ Planned Architecture (Blueprint)
*As specified in `PRD.md`:*
- **Architecture**: Microservices (Auth, Store, Inventory, Procurement, Sales, Finance).
- **Backend**: Go (Golang) with Echo Framework.
- **Persistence**: PostgreSQL (Isolation per service/database).
- **Caching & Messaging**: Redis (Streams for Event-Driven).
- **API Gateway**: Traefik with path-based routing.
- **Deployment**: Docker Swarm with Multi-stage builds.

## 🛠️ Implementation Pattern (Standard)
Setiap microservice baru harus mengikuti struktur folder dan layer berikut:

1.  **Directory Structure**:
    - `cmd/api/main.go`: Entrypoint, setup Echo, DB connection, and DI.
    - `internal/domain/`: Struct models dan Repository interfaces.
    - `internal/repository/`: Implementasi GORM repository.
    - `internal/service/`: Business logic, validation, dan unit tests.
    - `internal/handler/`: Echo handlers (REST endpoints).

2.  **Implementation Steps**:
    - **Step 1: Domain & Interface**. Definisikan struct dan interface di `internal/domain`.
    - **Step 2: Repository**. Implementasi database operations menggunakan GORM.
    - **Step 3: Service & Logic**. Implementasi business logic dan wajib menyertakan unit test (`*_test.go`).
    - **Step 4: Handler**. Implementasi REST API menggunakan framework Echo.
    - **Step 5: Main App & Containerization**. Setup `main.go`, `Dockerfile`, dan update `docker-compose.yml`.

3.  **Database Isolation**:
    - Gunakan database terpisah untuk setiap service (contoh: `auth_db`, `store_db`).
    - Gunakan script `docker/postgres/init-db.sh` untuk inisialisasi database di Docker.

4.  **Routing (Traefik)**:
    - Gunakan path-based routing di `docker-compose.yml` (contoh: `/auth`, `/stores`).
    - Pastikan menggunakan middleware `stripprefix` agar service menerima path yang bersih.

## 📜 Mandatory Instructions
- **BACA GEMINI.md**: Selalu baca file ini sebelum membuat Plan atau mulai Implementasi fitur/service baru.
- **Ikuti Pola**: Jangan menyimpang dari pola implementasi yang sudah ditetapkan di atas untuk menjaga konsistensi codebase.
- **Automated Verification**: Setiap service baru wajib memiliki unit test pada layer service untuk memvalidasi business logic utama.
- **Git Workflow**: Gunakan `git worktree` untuk setiap fitur baru guna menjaga kebersihan directory utama.

---
*Note: This GEMINI.md is the source of truth for implementation standards. Update it as the architecture evolves.*
