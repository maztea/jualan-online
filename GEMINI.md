# Project Context: Jualan Online Plan

This directory serves as the **Planning and Architectural Design** phase for the "Sistem Manajemen Inventaris & Finansial Multi-Channel". It contains the foundational requirements and microservices architecture specifications.

## 📁 Directory Overview
This project is currently in the **discovery and design stage**. It holds research notes, business requirements (PRD), and architectural blueprints for a system designed to centralize inventory and financial tracking across multiple online marketplaces (Shopee, Tokopedia, TikTok Shop, etc.).

## 🔑 Key Files
- **`PRD.md`**: The core Product Requirement Document. It defines the problem statement, user personas, functional requirements (user stories), and technical specifications (Microservices decomposition, Tech Stack, and Data Contracts).
- **`master_prompt.txt`**: Contains the strategic persona and structural guidelines used to generate the enterprise-grade PRD.
- **`.discovery/notes.txt`**: The raw initial input and business problems that triggered the project.

## 🏗️ Planned Architecture (Blueprint)
*As specified in `PRD.md`:*
- **Architecture**: Microservices (Auth, Store, Inventory, Procurement, Sales, Finance).
- **Backend**: Go (Golang).
- **Persistence**: PostgreSQL (Isolation per service).
- **Caching & Messaging**: Redis (Streams for Event-Driven).
- **API Gateway**: Traefik.
- **Deployment**: Docker Swarm with Multi-stage builds.

## 🛠️ Usage & Future Development
1.  **Reference the PRD**: Use `PRD.md` as the "Source of Truth" for any implementation tasks.
2.  **Implementation Phase**: 
    - **Git Workflow**: Selalu gunakan `git worktree` untuk setiap memulai implementasi fitur baru guna menjaga kebersihan *working directory* utama.
    - Create subdirectories for each microservice (e.g., `/services/auth-service`, `/services/inventory-service`) following the Tech Stack specified in the PRD.
3.  **CI/CD**: Follow the deployment strategy outlined in Section 6 of the PRD for automated linting, testing, and rolling updates.

---
*Note: This GEMINI.md is intended for instructional context and should be updated as the project transitions from planning to implementation.*
