# Split Repositories - Quick Start Guide

## 🎯 Goal

Split this monorepo into two separate repositories:
1. **Backend**: `osrs-price-api` (Go API)
2. **Frontend**: `osrs-ge-tracker` (Next.js app)

Each with its own Git history and CI/CD pipeline.

## 📋 Prerequisites

- Git installed
- GitHub account
- GitHub CLI (optional): `brew install gh`

## 🚀 Quick Setup (Automated)

I've created setup scripts for you!

### Option A: Use Setup Scripts (Easiest)

#### Backend:
```bash
cd /Users/ori/GolandProjects/Project1
./setup-backend-repo.sh
```

#### Frontend:
```bash
cd /Users/ori/GolandProjects/Project1/osrs-ge-tracker
./setup-frontend-repo.sh
```

Follow the prompts and you're done! 🎉

---

### Option B: Manual Setup (Step-by-step)

## 1️⃣ Set Up Backend Repository

```bash
# Navigate to backend root
cd /Users/ori/GolandProjects/Project1

# Initialize Git
git init

# Stage backend files (excluding frontend)
git add .gitignore .gitattributes .github/
git add main.go go.mod go.sum
git add internal/ cmd/ migrations/
git add *.md .env.example Makefile docker-compose.yml

# Commit
git commit -m "Initial commit: Go backend with CI/CD"

# Create GitHub repo (via web or CLI)
# Web: https://github.com/new → Create "osrs-price-api"
# CLI: gh repo create YOUR-USERNAME/osrs-price-api --public --source=. --push

# Link and push
git remote add origin https://github.com/YOUR-USERNAME/osrs-price-api.git
git branch -M main
git push -u origin main
```

## 2️⃣ Set Up Frontend Repository

```bash
# Navigate to frontend
cd /Users/ori/GolandProjects/Project1/osrs-ge-tracker

# Remove old git if exists
rm -rf .git

# Initialize new Git
git init

# Add all frontend files
git add .

# Commit
git commit -m "Initial commit: Next.js frontend"

# Create GitHub repo
# Web: https://github.com/new → Create "osrs-ge-tracker"
# CLI: gh repo create YOUR-USERNAME/osrs-ge-tracker --public --source=. --push

# Link and push
git remote add origin https://github.com/YOUR-USERNAME/osrs-ge-tracker.git
git branch -M main
git push -u origin main
```

## 3️⃣ Verify Everything Works

### Backend:
```bash
# Clone your new backend repo
git clone https://github.com/YOUR-USERNAME/osrs-price-api.git
cd osrs-price-api

# Test build
go build -o osrs-price-api main.go
./osrs-price-api
```

### Frontend:
```bash
# Clone your new frontend repo
git clone https://github.com/YOUR-USERNAME/osrs-ge-tracker.git
cd osrs-ge-tracker

# Test run
npm install
npm run dev
```

## 4️⃣ CI/CD Pipeline (Backend)

Once pushed to GitHub:

1. **Go to your backend repo** → Actions tab
2. **You should see**: "Build Go Backend" workflow
3. **On every push to main**: Automatic builds for:
   - ✅ Linux (amd64, arm64)
   - ✅ macOS (Intel, Apple Silicon)
   - ✅ Windows (amd64)

### Download Builds:

```bash
# Via web
https://github.com/YOUR-USERNAME/osrs-price-api/releases

# Via GitHub CLI
gh release download --repo YOUR-USERNAME/osrs-price-api

# Direct download (replace USERNAME and TAG)
curl -L https://github.com/YOUR-USERNAME/osrs-price-api/releases/latest/download/osrs-price-api-linux-amd64 -o osrs-api
chmod +x osrs-api
./osrs-api
```

## 📁 New Directory Structure

### Before:
```
Project1/
├── (backend files)
└── osrs-ge-tracker/
    └── (frontend files)
```

### After:
```
Your GitHub Account/
├── osrs-price-api/         (Backend repo)
│   ├── .github/workflows/  (CI/CD)
│   ├── internal/
│   ├── migrations/
│   ├── main.go
│   └── ...
└── osrs-ge-tracker/        (Frontend repo)
    ├── src/
    ├── public/
    ├── package.json
    └── ...
```

## 🔧 Local Development Setup

After separation, organize your local workspace:

```bash
mkdir ~/Projects/osrs-project
cd ~/Projects/osrs-project

# Clone both repos
git clone https://github.com/YOUR-USERNAME/osrs-price-api.git backend
git clone https://github.com/YOUR-USERNAME/osrs-ge-tracker.git frontend

# Start backend
cd backend
go run main.go

# In another terminal, start frontend
cd ../frontend
npm run dev
```

## 🎯 Benefits

✅ **Independent versioning** - Each repo has its own version tags  
✅ **Separate CI/CD** - Backend builds executables, frontend deploys to Vercel  
✅ **Better collaboration** - Different teams can work independently  
✅ **Cleaner history** - Each repo has relevant commits only  
✅ **Smaller clones** - Faster git operations  
✅ **Platform-specific** - Backend can deploy anywhere, frontend on Vercel  

## 🚨 Important Notes

### Environment Variables

**Backend** (`.env`):
```env
DATABASE_URL=postgresql://...
PORT=8080
```

**Frontend** (`.env.local`):
```env
NEXT_PUBLIC_GO_API_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=...
NEXT_PUBLIC_SUPABASE_ANON_KEY=...
```

### Update Frontend Configuration

After deploying backend, update frontend:
```env
# Production
NEXT_PUBLIC_GO_API_URL=https://your-backend.railway.app

# Development
NEXT_PUBLIC_GO_API_URL=http://localhost:8080
```

## 📦 GitHub Releases (Backend)

Every push to `main` creates a release with binaries:

```
osrs-price-api/releases/
├── osrs-price-api-linux-amd64      (Linux)
├── osrs-price-api-linux-arm64      (Linux ARM)
├── osrs-price-api-darwin-amd64     (macOS Intel)
├── osrs-price-api-darwin-arm64     (macOS Apple Silicon)
└── osrs-price-api-windows-amd64.exe (Windows)
```

### Usage:
```bash
# Download
wget https://github.com/YOUR-USERNAME/osrs-price-api/releases/latest/download/osrs-price-api-linux-amd64

# Run
chmod +x osrs-price-api-linux-amd64
./osrs-price-api-linux-amd64
```

## 🔄 Development Workflow

### Backend Changes:
```bash
cd backend
# Make changes
git add .
git commit -m "Add new feature"
git push
# → GitHub Actions builds new binaries automatically
```

### Frontend Changes:
```bash
cd frontend
# Make changes
git add .
git commit -m "Update UI"
git push
# → Vercel deploys automatically (if connected)
```

## 🐛 Troubleshooting

### "fatal: not a git repository"
```bash
git init
```

### GitHub Actions not showing up
1. Check `.github/workflows/build.yml` exists in backend
2. Go to repo Settings → Actions → Enable workflows
3. Push a commit to trigger

### Frontend can't connect to backend
1. Check `NEXT_PUBLIC_GO_API_URL` in `.env.local`
2. Ensure backend is running: `curl http://localhost:8080/health`
3. Check CORS is enabled in backend (already done)

## ✅ Checklist

**Backend Setup**:
- [ ] Git initialized
- [ ] Files committed
- [ ] Pushed to GitHub
- [ ] GitHub Actions working
- [ ] First build completed
- [ ] Binary downloaded and tested

**Frontend Setup**:
- [ ] Git initialized
- [ ] Files committed
- [ ] Pushed to GitHub
- [ ] Environment variables set
- [ ] Backend URL configured
- [ ] App runs successfully

## 🎓 Next Steps

1. ✅ Set up CI/CD for frontend (Vercel auto-deploy)
2. ✅ Deploy backend to production (Railway/Fly.io)
3. ✅ Configure production environment variables
4. ✅ Set up monitoring (Sentry, etc.)
5. ✅ Add API documentation
6. ✅ Set up automatic backups

## 📚 Documentation

- **Backend**: See `README.md` in backend repo
- **Frontend**: See `README.md` in frontend repo
- **API Docs**: See backend `/docs` folder
- **Full Setup**: See `FULL_STACK_SETUP.md`

---

Need help? Check `REPO_SEPARATION_GUIDE.md` for detailed instructions!