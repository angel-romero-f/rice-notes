# CLAUDE.md
**Guidance for Claude Code (claude.ai/code)**  
Always load this context before generating or refactoring code for this repository.

---

## 1 · Project Overview
**Rice Notes** is a full-stack note-sharing platform for Rice University students.  
Flow → Google sign-in (`@rice.edu` only) → upload / view course-tagged notes (PDF, images, markdown) → future AI summarization & search.

---

## 2 · Tech Stack

| Layer / Concern | Technology & Version | Notes |
|-----------------|----------------------|-------|
| **Backend**     | Go 1.22 · chi v5     | idiomatic Go |
| Auth            | Google OAuth 2 → self-issued **JWT** in HttpOnly cookie | reject non-`@rice.edu` emails |
| Storage         | AWS S3 (SDK v2) for objects · Postgres 15 via pgxpool | local dev → MinIO + Postgres (docker-compose) |
| Front-end       | Next.js 15.3.1 (React 19) · TypeScript · Tailwind CSS v4 | deployed on Vercel |
| DevOps / CI     | Dockerfile · docker-compose · GitHub Actions | CI runs `go test ./...` & `npm run build` |

---

## 3 · Repository Layout

backend/
├ cmd/server/main.go # starts HTTP server (port 3000)
└ internal/
├ handlers/ # HTTP controllers
├ services/ # business logic + interfaces
├ repository/ # Postgres queries (TODO)
├ infra/storage/ # S3 adapter implementing Uploader (TODO)
├ middleware/ # CORS, JWT, RequestID
├ routes/ # chi router & dependency wiring
└ models/ # domain structs
frontend/ # Next.js app directory
migrations/ # .sql files (golang-migrate)
pkg/logger/ # PrettyHandler (slog)



> **Dependency-injection rule:** `routes` wires **concrete repo → service → handler**.  
> Interfaces live in the *consumer* package; lower layers never import higher ones.

---

## 4 · Current State
`GET /` returns a welcome message via `NoteHandler`.  
Repository layer, storage adapters, and full auth flow are still **TODO**.

---

## 5 · MVP Road-map (vertical slices)

1. **Auth slice**  
   • `/api/auth/google` → redirect  
   • `/api/auth/google/callback` → code exchange, `@rice.edu` check, issue JWT cookie  
2. **Note upload**  
   • `POST /api/notes` (multipart) → upload to S3 key `notes/<user>/<uuid>`; insert metadata row  
3. **List & download**  
   • `GET /api/notes?course_id=` → list notes  
   • `GET /api/notes/{id}` → 302 to presigned S3 URL  
4. **Future** → OpenAI summarization & pgvector search

---

## 6 · Coding Guidelines

* All service / repository funcs: **first arg** `ctx context.Context`.  
* No circular imports; lower layers **never** import higher ones.  
* Use raw pgx SQL (`$1,$2…`), no ORM.  
* Handlers map service errors → HTTP status; services return plain errors.  
* Exported symbols must have GoDoc.  
* Unit tests: `httptest` + mocks; integration: Dockertest / Testcontainers.  
* Keep PR-sized output (< 150 LOC) unless user requests more.

---

## 7 · Dev Commands

### Front-end
```bash
cd frontend
npm run dev       # dev server
npm run build     # production build
npm run start     # run prod build
npm run lint
Back-end
bash
Copy
Edit
cd backend
go run cmd/server/main.go          # dev server
go build -o rice-notes ./cmd/...   # build binary
go test ./...                      # tests
docker-compose up                  # Postgres + MinIO for local dev