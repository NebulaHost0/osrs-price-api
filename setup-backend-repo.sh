#!/bin/bash
# Script to set up the backend repository

set -e

echo "🚀 Setting up OSRS Price API Backend Repository"
echo "================================================"
echo ""

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo "⚠️  GitHub CLI not found. Install it with: brew install gh"
    echo "Or create the repository manually on GitHub.com"
    CREATE_REPO_MANUALLY=true
else
    CREATE_REPO_MANUALLY=false
fi

# Get GitHub username
read -p "Enter your GitHub username: " GITHUB_USER

# Repository name
REPO_NAME="osrs-price-api"

echo ""
echo "📦 Repository will be: https://github.com/$GITHUB_USER/$REPO_NAME"
echo ""

# Initialize git if not already
if [ ! -d .git ]; then
    echo "📝 Initializing Git repository..."
    git init
else
    echo "✅ Git already initialized"
fi

# Create .gitignore for backend
echo "📝 Creating .gitignore..."
cat > .gitignore << 'EOF'
# Binaries
*.exe
*.dll
*.so
*.dylib
bin/
dist/
osrs-price-api
osrs-price-api-*

# Test files
*.test
*.out
coverage.txt

# Dependencies
vendor/

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
.DS_Store

# Logs
*.log
logs/

# Temporary
tmp/
temp/

# Exclude frontend
osrs-ge-tracker/
EOF

# Create .gitattributes
cat > .gitattributes << 'EOF'
# Go files
*.go text eol=lf

# Scripts
*.sh text eol=lf

# SQL migrations
*.sql text eol=lf

# Config files
*.yaml text eol=lf
*.yml text eol=lf
*.json text eol=lf
*.toml text eol=lf

# Markdown
*.md text eol=lf
EOF

echo "📝 Adding backend files to git..."
git add .gitignore .gitattributes
git add .github/
git add main.go go.mod go.sum
git add internal/ cmd/ migrations/
git add *.md
git add .env.example Makefile docker-compose.yml

# Create README.md from BACKEND_README.md
if [ -f BACKEND_README.md ]; then
    cp BACKEND_README.md README.md
    git add README.md
fi

# Check for changes
if git diff --cached --quiet; then
    echo "⚠️  No changes to commit"
else
    echo "💾 Committing backend files..."
    git commit -m "Initial commit: OSRS Price API backend with CI/CD

- Go API with Gin framework
- PostgreSQL database with GORM
- Automatic price collection from OSRS Wiki
- Volume tracking and trade analysis
- Smart data aggregation (saves 95% on costs)
- GitHub Actions CI/CD pipeline
- Multi-platform binary builds"
fi

# Create GitHub repository
if [ "$CREATE_REPO_MANUALLY" = false ]; then
    echo ""
    read -p "Create GitHub repository now? (y/n): " CREATE_NOW
    if [ "$CREATE_NOW" = "y" ]; then
        echo "🔨 Creating GitHub repository..."
        gh repo create $GITHUB_USER/$REPO_NAME --public --source=. --remote=origin --push
        echo "✅ Repository created and pushed!"
    else
        echo "ℹ️  Skipping repository creation. Create it manually, then run:"
        echo "   git remote add origin https://github.com/$GITHUB_USER/$REPO_NAME.git"
        echo "   git push -u origin main"
    fi
else
    echo ""
    echo "📋 Next steps:"
    echo "1. Create repository on GitHub: https://github.com/new"
    echo "2. Name it: $REPO_NAME"
    echo "3. Run these commands:"
    echo "   git remote add origin https://github.com/$GITHUB_USER/$REPO_NAME.git"
    echo "   git push -u origin main"
fi

echo ""
echo "✅ Backend repository setup complete!"
echo "📁 Location: $(pwd)"
echo "🌐 GitHub: https://github.com/$GITHUB_USER/$REPO_NAME"
echo ""
echo "⚙️  GitHub Actions will automatically build binaries on every push to main"
echo "📦 Download builds from: https://github.com/$GITHUB_USER/$REPO_NAME/releases"