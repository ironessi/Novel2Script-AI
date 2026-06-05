# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Novel2Script-AI converts novels into structured YAML screenplays using a multi-stage AI pipeline: chapter splitting → character extraction (Character Bible) → plot event chain → scene planning → YAML script generation → schema validation → hallucination detection → safety review.

Design doc (Chinese): `Novel2Script-AI-项目设计文档.md`

## Architecture

Three-service system: GoFrame backend (port 8000) → Python FastAPI AI service (port 9000) → MySQL/Redis. Frontend (Vue3, port 5173) planned but not yet built.

```
Frontend → GoFrame Backend → Python AI Service → LLM (DeepSeek/Ollama)
               ↓                    ↓
             MySQL               Redis
          (business data)     (task state/progress)
```

## Backend Architecture (GoFrame v2)

Layered design — follow this pattern when adding new features:

- **`api/v1/`** — Request/response structs with GoFrame validation tags (`v:"required|..."`)
- **`controller/`** — Thin HTTP handlers. Parse request, call service, write JSON response. Access path params via `r.GetRouter("id")`.
- **`service/`** — Interface definitions + package-level singleton vars (`var Auth IAuth`)
- **`logic/`** — Business logic implementing service interfaces. Register via `init()` setting the service var.
- **`dao/`** — Database access using GoFrame ORM (`db.Model("table").Ctx(ctx)...`). Use `One()` for single records (returns nil on empty), `Scan(&slice)` for lists.
- **`model/entity/`** — Database entity structs

**Adding a new feature:** Create API types → entity model → DAO functions → service interface → logic implementation → controller → register route in `cmd.go`.

**Service registration pattern:**
```go
// In logic/xxx/xxx.go
func init() { service.Xxx = &xxxImpl{} }
```
Then blank-import in `cmd.go`: `_ "novel2script-backend/internal/logic/xxx"`

## Commands

```bash
# Backend
cd backend && go run main.go          # Start server (reads .env from project root)
cd backend && go build ./...          # Check compilation
cd backend && go mod tidy             # Sync dependencies

# Infrastructure
cd deploy && docker compose up -d mysql redis   # Start MySQL + Redis

# Database
mysql -u root -p < deploy/mysql-init.sql        # Schema (10 tables)
mysql -u root -p novel2script < scripts/seed.sql # Test data

# AI Service (not yet implemented)
cd ai-service && uvicorn app.main:app --host 0.0.0.0 --port 9000 --reload

# Frontend (not yet implemented)
cd frontend && npm run dev
```

## Key Configuration

- `.env` at project root — MySQL credentials, JWT secret, Redis, AI service URL/token
- Backend loads `.env` automatically via `config.Init()` (searches `../../.env`, `../.env`, `.env`)
- Redis uses DB 2 (not default 0)
- JWT tokens expire in 24 hours, HS256 signing

## Database

- 10 tables, all BIGINT auto-increment PKs, `utf8mb4_unicode_ci`
- `novel_project` uses soft delete (`deleted_at` column) — always filter `WHERE deleted_at IS NULL`
- JSON columns in `character_profile` (aliases, personality, relationships, source_refs) and `plot_event` (source_refs) use custom `entity.JSON` type
- All project queries must include `owner_id` check to prevent unauthorized access

## Security Rules

- All SQL via GoFrame ORM or parameterized queries — never string concatenation
- Project access: `owner_id == current_user OR role == "admin"` (enforced in logic layer)
- Passwords: bcrypt only, never store plaintext
- File uploads (when implemented): limit to `.txt/.md/.docx`, max 20MB, UUID filenames, path traversal prevention
- JWT secret and DB passwords from `.env` only, never hardcoded

## Incomplete / Placeholder

- `ai-service/` — not created yet (Python FastAPI)
- `frontend/` — not created yet (Vue3)
- `utility/filecheck/`, `utility/sanitizer/` — empty directories
- Task creation (`task.Create`) only writes to DB, does not trigger AI processing
- Script validate/export are stubs (no AI service integration)
- File upload endpoint not implemented
- No tests exist yet

## Ports

| Service | Port |
|---------|------|
| Backend | 8000 |
| AI Service | 9000 |
| Frontend | 5173 |
| MySQL | 3306 |
| Redis | 6379 |
