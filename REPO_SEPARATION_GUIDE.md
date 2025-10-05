# Repository Separation Guide

This guide will help you split the backend and frontend into separate Git repositories.

## Overview

**Current structure**:
```
Project1/
├── (Go backend files)
└── osrs-ge-tracker/ (Next.js frontend)
```

**Target structure**:
- Repository 1: `osrs-price-api` (Backend)
- Repository 2: `osrs-ge-tracker` (Frontend)

## Step-by-Step Instructions

### Part 1: Prepare Backend Repository

1. **Create a new GitHub repository** for the backend:
   - Go to GitHub → New Repository
   - Name: `osrs-price-api`
   - Description: "Go backend API for OSRS Grand Exchange price tracking"
   - Public or Private (your choice)
   - **Do NOT** initialize with README, .gitignore, or license

2. **Navigate to your backend directory**:
```bash
cd /Users/ori/GolandProjects/Project1
```

3. **Initialize Git** (if not already):
```bash
git init
```

4. **Add backend files** to Git (excluding frontend):
```bash
# Add .gitignore first
git add .gitignore

# Add GitHub Actions workflow
git add .github/

# Add all backend files
git add main.go go.mod go.sum
git add internal/
git add cmd/
git add migrations/
git add *.md
git add .env.example
git add Makefile
git add docker-compose.yml

# Commit
git commit -m "Initial backend setup with CI/CD pipeline"
```

5. **Link to remote and push**:
```bash
# Replace with your actual GitHub username
git remote add origin https://github.com/YOUR-USERNAME/osrs-price-api.git
git branch -M main
git push -u origin main
```

### Part 2: Prepare Frontend Repository

1. **Create another GitHub repository** for frontend:
   - Name: `osrs-ge-tracker`
   - Description: "Next.js frontend for OSRS Grand Exchange tracking"
   - Public or Private
   - **Do NOT** initialize with README

2. **Navigate to frontend directory**:
```bash
cd /Users/ori/GolandProjects/Project1/osrs-ge-tracker
```

3. **Initialize Git for frontend**:
```bash
# Check if .git already exists
ls -la

# If .git exists and you want fresh history, remove it:
rm -rf .git

# Initialize new repo
git init

# Create frontend-specific .gitignore
cat > .gitignore << 'EOF'
# Dependencies
node_modules/
.pnp
.pnp.js

# Testing
coverage/

# Next.js
.next/
out/
build/
dist/

# Production
*.tsbuildinfo

# Misc
.DS_Store
*.pem

# Debug
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Local env files
.env
.env*.local

# Vercel
.vercel

# Typescript
*.tsbuildinfo
next-env.d.ts
EOF

# Add all frontend files
git add .

# Commit
git commit -m "Initial frontend setup with Go API integration"
```

4. **Link to remote and push**:
```bash
git remote add origin https://github.com/YOUR-USERNAME/osrs-ge-tracker.git
git branch -M main
git push -u origin main
```

### Part 3: Update Configuration

#### Backend `.env.example`:
```env
# Database Configuration
DATABASE_URL=postgresql://postgres:password@localhost:5432/osrs_prices?sslmode=disable

# Server Configuration
PORT=8080
```

#### Frontend `.env.local.example`:
```env
# Go API Configuration
NEXT_PUBLIC_GO_API_URL=http://localhost:8080

# Supabase Configuration (for auth)
NEXT_PUBLIC_SUPABASE_URL=your-supabase-url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
```

### Part 4: Clean Up Original Directory

After successfully pushing both repos:

```bash
# Back up original directory
cd /Users/ori/GolandProjects
mv Project1 Project1-backup

# Clone the new repos
mkdir Project1
cd Project1

git clone https://github.com/YOUR-USERNAME/osrs-price-api.git backend
git clone https://github.com/YOUR-USERNAME/osrs-ge-tracker.git frontend

# Now you have:
# Project1/backend/ (Go API)
# Project1/frontend/ (Next.js app)
```

### Part 5: Update README Files

#### Backend `README.md`:
Already created as `BACKEND_README.md` - rename it:
```bash
cd backend
mv BACKEND_README.md README.md
git add README.md
git commit -m "Add comprehensive README"
git push
```

#### Frontend `README.md`:
Update the frontend README to point to the backend repo:
```markdown
# OSRS Grand Exchange Tracker - Frontend

Next.js frontend for the OSRS Grand Exchange price tracker.

## Backend

This frontend requires the Go backend API running. Get it here:
👉 [osrs-price-api](https://github.com/YOUR-USERNAME/osrs-price-api)

## Setup

1. Clone the repository
2. Install dependencies: `npm install`
3. Copy `.env.local.example` to `.env.local`
4. Update `NEXT_PUBLIC_GO_API_URL` to your backend URL
5. Run: `npm run dev`

See full documentation in the repository.
```

## CI/CD Pipeline (Backend)

The backend now has automatic builds via GitHub Actions:

### What It Does:
1. **On every push to `main`**:
   - Runs tests
   - Builds executables for:
     - Linux (amd64, arm64)
     - macOS (Intel, Apple Silicon)
     - Windows (amd64)
   - Creates a GitHub Release with all binaries

2. **To download builds**:
```bash
# Visit your repo releases page:
https://github.com/YOUR-USERNAME/osrs-price-api/releases

# Or use GitHub CLI:
gh release download --repo YOUR-USERNAME/osrs-price-api

# Or curl:
curl -L https://github.com/YOUR-USERNAME/osrs-price-api/releases/latest/download/osrs-price-api-linux-amd64 -o osrs-price-api
chmod +x osrs-price-api
./osrs-price-api
```

## Verifying Everything Works

### Backend:
```bash
cd backend
go run main.go
# Should start on http://localhost:8080
```

### Frontend:
```bash
cd frontend
npm run dev
# Should start on http://localhost:3000
```

### Test the connection:
```bash
# Backend health check
curl http://localhost:8080/health

# Frontend should load and connect to backend
open http://localhost:3000
```

## Benefits of Separated Repos

✅ **Independent versioning** - Backend and frontend have separate release cycles
✅ **Cleaner CI/CD** - Each repo has its own build pipeline
✅ **Better access control** - Different team permissions per repo
✅ **Smaller repo sizes** - Faster cloning and operations
✅ **Specialized .gitignore** - Each repo ignores only what it needs
✅ **Independent deployment** - Deploy backend/frontend separately

## Troubleshooting

### Backend repo too large?
```bash
# Remove large files from git history
git filter-branch --tree-filter 'rm -rf osrs-ge-tracker' HEAD
git push origin main --force
```

### Frontend missing files?
```bash
# Make sure you're in the right directory
pwd  # Should show .../osrs-ge-tracker

# Check what's being tracked
git ls-files
```

### GitHub Actions not running?
1. Check `.github/workflows/build.yml` exists in backend repo
2. Go to repo → Actions tab → Enable workflows
3. Push a new commit to trigger build

## Next Steps

1. ✅ Set up backend repo with CI/CD
2. ✅ Set up frontend repo
3. 📝 Update documentation links between repos
4. 🚀 Deploy backend (Railway/DigitalOcean/AWS)
5. 🌐 Deploy frontend (Vercel/Netlify)
6. 🔐 Set up environment variables in deployment platforms
7. 📊 Monitor builds and deployments

## Support

Questions? Open an issue in the respective repository:
- Backend issues: `osrs-price-api/issues`
- Frontend issues: `osrs-ge-tracker/issues`