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
| Front-end       | Next.js 15.3.1 (React 19) · TypeScript · Tailwind CSS v4 | deployed on Vercel: https://rice-notes.vercel.app |
| DevOps / CI     | Dockerfile · docker-compose · GitHub Actions | CI runs `go test ./...` & `npm run build` |

---

## 3 · Repository Layout

backend/
├ cmd/server/main.go # starts HTTP server (port 8080)
└ internal/
├ handlers/ # HTTP controllers (note_handler.go, auth_handler.go)
├ services/ # business logic + interfaces
├ repository/ # Postgres queries (note_repository.go) 
├ infra/storage/ # S3 adapter implementing Uploader (s3_uploader.go)
├ middleware/ # CORS, JWT, RequestID
├ routes/ # chi router & dependency wiring
└ models/ # domain structs (note.go)
frontend/ # Next.js app directory
├ src/app/ # App router pages (dashboard, etc.)
├ src/components/ # UI components (UploadModal, CourseSelect, etc.)
├ src/hooks/ # React hooks (useAuth, useFileUpload)
└ src/lib/ # API service layer (api.ts)
migrations/ # .sql files (001_create_notes_table.sql)
pkg/logger/ # PrettyHandler (slog)



> **Dependency-injection rule:** `routes` wires **concrete repo → service → handler**.  
> Interfaces live in the *consumer* package; lower layers never import higher ones.

---

## 4 · Current State
**✅ Complete PDF Upload System** implemented with:
- **Auth**: Google OAuth 2.0 with JWT in HttpOnly cookies (`@rice.edu` validation)
- **Backend**: Complete CRUD API with PostgreSQL + S3 storage
- **Frontend**: Drag & drop upload modal with progress tracking
- **Database**: Notes table with proper indexing and relationships
- **Storage**: AWS S3 integration with presigned URLs and mock support

**API Endpoints**:
- `POST /api/notes` - PDF upload (multipart form-data)
- `GET /api/notes` - List user's notes (with pagination & filtering) 
- `GET /api/notes/{id}` - Get specific note with presigned download URL
- `DELETE /api/notes/{id}` - Delete note (removes from S3 + database)

---

## 5 · MVP Road-map (vertical slices)

1. **✅ Auth slice** (COMPLETED)  
   • `/api/auth/google` → redirect to Google OAuth  
   • `/api/auth/google/callback` → code exchange, `@rice.edu` check, issue JWT cookie  
   • `/api/auth/me` → get user info from JWT  
   
2. **✅ Note upload** (COMPLETED)  
   • `POST /api/notes` (multipart) → upload PDF to S3, store metadata in PostgreSQL  
   • Frontend drag & drop interface with progress tracking  
   • File validation (PDF only, 10MB limit) & course selection  
   
3. **✅ List & download** (COMPLETED)  
   • `GET /api/notes?course_id=&limit=&offset=` → paginated list with filtering  
   • `GET /api/notes/{id}` → presigned S3 download URL  
   • `DELETE /api/notes/{id}` → remove from S3 + database  
   
4. **🔄 Future Enhancements**  
   • OpenAI summarization & content extraction  
   • pgvector search with semantic similarity  
   • Note sharing between users  
   • Comment system and ratings

---

## 6 · Database Schema

**Notes Table** (`notes`):
```sql
- id: UUID PRIMARY KEY
- user_email: VARCHAR(255) NOT NULL (from JWT claims)
- title: VARCHAR(255) NOT NULL  
- course_id: VARCHAR(20) NOT NULL (e.g., "COMP101", "MATH220")
- file_name: VARCHAR(255) NOT NULL (original filename)
- file_path: TEXT NOT NULL (S3 key: "notes/{email}/{uuid}/{filename}")
- file_size: BIGINT NOT NULL (bytes)
- content_type: VARCHAR(50) NOT NULL (application/pdf)
- uploaded_at: TIMESTAMPTZ DEFAULT NOW()
- updated_at: TIMESTAMPTZ DEFAULT NOW()

Indexes: user_email, course_id, uploaded_at
```

**Environment Variables** (`backend/.local.env`):
```bash
# Database - Docker PostgreSQL
DATABASE_URL=postgres://postgres:password@localhost:5432/rice_notes?sslmode=disable

# AWS S3 - Production Setup
AWS_REGION=us-east-2
AWS_S3_BUCKET=rice-notes-prod
AWS_ACCESS_KEY_ID=your_actual_access_key
AWS_SECRET_ACCESS_KEY=your_actual_secret_key
USE_MOCK_S3=false

# Google OAuth - Your existing credentials
GOOGLE_CLIENT_ID="your_google_client_id"
GOOGLE_CLIENT_SECRET="your_google_client_secret"
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# JWT Secret - Your existing secret
JWT_SECRET=your_jwt_secret_key

# Development settings
ENV=development
```

---

## 7 · Coding Guidelines

* All service / repository funcs: **first arg** `ctx context.Context`.  
* No circular imports; lower layers **never** import higher ones.  
* Use raw pgx SQL (`$1,$2…`), no ORM.  
* Handlers map service errors → HTTP status; services return plain errors.  
* Exported symbols must have GoDoc.  
* Unit tests: `httptest` + mocks; integration: Dockertest / Testcontainers.  
* Keep PR-sized output (< 150 LOC) unless user requests more.
* Always test the frontend cod by using the playwright MCP to test the desired workflow of the feature. 

---

## 8 · Local Development Setup

### Prerequisites
```bash
# Install required tools
brew install golang-migrate docker
```

### Database Setup (PostgreSQL in Docker)
```bash
# Start PostgreSQL container
docker run --name rice-notes-postgres \
  -e POSTGRES_DB=rice_notes \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15

# Check container is running
docker ps

# Run database migration
cd backend
migrate -database "postgres://postgres:password@localhost:5432/rice_notes?sslmode=disable" -path migrations up

# Verify tables created
docker exec -it rice-notes-postgres psql -U postgres -d rice_notes -c "\dt"
```

### AWS S3 Setup
1. **Create S3 bucket**: `rice-notes-prod` in `us-east-2`
2. **Create IAM user**: `rice-notes-app` with S3 permissions
3. **Add credentials** to `backend/.local.env`

### Start Development Servers
```bash
# Terminal 1 - Backend (port 8080)
cd backend
go run cmd/server/main.go

# Terminal 2 - Frontend (port 3000)
cd frontend
npm run dev
```

### Test Local Setup
1. Visit `http://localhost:3000`
2. Sign in with `@rice.edu` Google account
3. Click "Upload Notes" in dashboard
4. Upload a PDF with title and course ID
5. Verify upload success and file appears in user's notes

---

## 9 · Production Deployment

### Docker Containers on EC2
```bash
# Backend API container
docker build -t rice-notes-api ./backend
docker run -p 8080:8080 --env-file .env rice-notes-api

# PostgreSQL container (same as local)
docker run --name rice-notes-postgres \
  -e POSTGRES_DB=rice_notes \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=production_password \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  -d postgres:15
```

### Environment Variables (Production)
- Same as local but update:
- `DATABASE_URL` → Use production password
- `GOOGLE_REDIRECT_URL` → Use production domain
- `AWS_REGION/BUCKET` → Production S3 bucket

---

## 10 · Dev Commands Reference

### Backend Development
```bash
cd backend

# Database operations
migrate -database $DATABASE_URL -path migrations up
migrate -database $DATABASE_URL -path migrations down 1

# Server operations
go run cmd/server/main.go          # dev server (:8080)
go build -o rice-notes ./cmd/...   # build binary
go test ./...                      # run tests
go mod tidy                        # update dependencies
```

### Frontend Development
```bash
cd frontend
npm run dev       # dev server (:3000)
npm run build     # production build
npm run start     # run prod build
npm run lint      # ESLint
```

### Docker Database Management
```bash
# Container management
docker start rice-notes-postgres
docker stop rice-notes-postgres
docker restart rice-notes-postgres

# Database access
docker exec -it rice-notes-postgres psql -U postgres -d rice_notes

# View logs
docker logs rice-notes-postgres
```

### Migration Management
```bash
# Create new migration
migrate create -ext sql -dir backend/migrations add_new_feature

# Check migration status
migrate -database $DATABASE_URL -path backend/migrations version

# Force migration version (if stuck)
migrate -database $DATABASE_URL -path backend/migrations force VERSION
```

---

## 11 · Production Deployment URLs

### Live Application URLs:
- **Frontend (Vercel)**: https://rice-notes.vercel.app
- **Backend API (EC2)**: https://ricenotesapi.com (HTTPS via Cloudflare)
- **Backend Direct**: http://3.141.195.234 (port 80, proxied by Cloudflare)
- **Database**: PostgreSQL container on EC2 (port 5432, internal access only)

### Production Environment:
- **EC2 Instance**: `3.141.195.234` (us-east-2)
- **Domain**: ricenotesapi.com (Namecheap → Cloudflare SSL)
- **Docker Containers**: `rice-notes-api-prod` & `rice-notes-postgres-prod`
- **Security Group**: Ports 22, 80, 443 open
- **Environment Variables**: Set in EC2 `/home/ec2-user/.bashrc`
- **Key File**: `~/.ssh/rice-notes-key2.pem`

### Production Access:
```bash
# SSH to production server
ssh -i ~/.ssh/rice-notes-key2.pem ec2-user@3.141.195.234

# Check container status
docker-compose -f docker-compose.prod.yml ps

# View API logs
docker-compose -f docker-compose.prod.yml logs api

# Restart services (IMPORTANT: Use -E flag to preserve environment variables)
sudo -E docker-compose -f docker-compose.prod.yml restart
```

### Production Deployment Process:
```bash
# 1. Upload updated files to EC2
scp -i ~/.ssh/rice-notes-key2.pem /path/to/file ec2-user@3.141.195.234:~/backend/internal/handlers/
scp -i ~/.ssh/rice-notes-key2.pem /path/to/docker-compose.prod.yml ec2-user@3.141.195.234:~/

# 2. SSH into EC2
ssh -i ~/.ssh/rice-notes-key2.pem ec2-user@3.141.195.234

# 3. Update environment variables if needed
echo 'export NEW_VAR="value"' >> ~/.bashrc
source ~/.bashrc

# 4. Rebuild and restart containers (CRITICAL: Use -E flag)
sudo -E docker-compose -f docker-compose.prod.yml down
sudo -E docker-compose -f docker-compose.prod.yml build api
sudo -E docker-compose -f docker-compose.prod.yml up -d
```

### Frontend Environment Configuration (Vercel):
The frontend is deployed on Vercel and requires environment variables to be set in the Vercel dashboard:

**Environment Variables to set in Vercel:**
- `NEXT_PUBLIC_API_URL=https://ricenotesapi.com` (Production API URL)

**Local Development:**
- Uses `.env.local` with `NEXT_PUBLIC_API_URL=http://localhost:8080`
- Never commit `.env.local` to git (already in .gitignore)
- Never create `.env.production` - use Vercel dashboard instead

### Important Notes:
- **ALWAYS use `sudo -E`** when running docker-compose commands to preserve environment variables
- Backend runs on **port 80** (not 8080) to work with Cloudflare SSL
- OAuth redirects to `FRONTEND_URL` environment variable (set to https://rice-notes.vercel.app)
- Google OAuth requires HTTPS - domain configured with Cloudflare SSL "Flexible" mode
- Files are uploaded to `~/backend/` directory structure on EC2 (not `~/rice-notes/backend/`)
- **Frontend environment variables are managed in Vercel dashboard, not in .env files**