# OSRS Price API - Backend

Go backend API for OSRS Grand Exchange price tracking with automatic data collection, volume tracking, and intelligent data aggregation.

## Features

- 🔄 Automatic price collection every 5 minutes from OSRS Wiki
- 📊 Real-time price data for 3,700+ items
- 📈 Historical price tracking with smart aggregation
- 💾 Volume tracking for trading activity analysis
- 🎯 Top gainers/losers and most traded items
- ⚡ Built-in caching for performance
- 🗄️ Automatic database maintenance (saves 95% on costs)
- 🌐 RESTful API with CORS support

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **External API**: OSRS Wiki Real-time Prices API

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**:
```bash
git clone <your-backend-repo-url>
cd osrs-price-api
```

2. **Set up environment variables**:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. **Install dependencies**:
```bash
go mod download
```

4. **Run migrations** (automatic on first start):
```bash
go run cmd/migrate/main.go
```

5. **Start the server**:
```bash
go run main.go
```

The API will start on `http://localhost:8080`

## Using Pre-built Binaries

Download the latest build from [GitHub Releases](../../releases):

```bash
# Example for Linux
wget https://github.com/your-username/osrs-price-api/releases/latest/download/osrs-price-api-linux-amd64
chmod +x osrs-price-api-linux-amd64
./osrs-price-api-linux-amd64

# Example for macOS (Apple Silicon)
wget https://github.com/your-username/osrs-price-api/releases/latest/download/osrs-price-api-darwin-arm64
chmod +x osrs-price-api-darwin-arm64
./osrs-price-api-darwin-arm64

# Example for Windows
# Download osrs-price-api-windows-amd64.exe and run it
```

## Environment Variables

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/osrs_prices?sslmode=disable

# Server
PORT=8080
```

## API Endpoints

### Current Prices
- `GET /api/v1/prices` - Get all current prices
- `GET /api/v1/prices/:id` - Get specific item price

### Historical Data
- `GET /api/v1/history/:id?hours=24` - Get price history
- `GET /api/v1/change/:id?hours=24` - Get price change
- `GET /api/v1/stats/:id?hours=168` - Get price statistics

### Market Analysis
- `GET /api/v1/gainers?limit=10&hours=24` - Top price gainers
- `GET /api/v1/losers?limit=10&hours=24` - Top price losers
- `GET /api/v1/volume?limit=10&hours=24` - Most traded items

### System
- `GET /health` - Health check
- `POST /api/v1/cache/clear` - Clear cache

## Database Maintenance

The system automatically:
- Keeps 7 days of 5-minute data (detailed)
- Aggregates to hourly for 8-90 days
- Aggregates to daily for 90+ days
- Maintains up to 5 years of history in ~1.7 GB

See [DATABASE_MAINTENANCE.md](DATABASE_MAINTENANCE.md) for details.

## Development

### Run Tests
```bash
go test -v ./...
```

### Build Binary
```bash
go build -o osrs-price-api main.go
```

### Build for Multiple Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/osrs-price-api-linux-amd64 main.go

# macOS
GOOS=darwin GOARCH=arm64 go build -o bin/osrs-price-api-darwin-arm64 main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/osrs-price-api-windows-amd64.exe main.go
```

## Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   ├── cache/             # In-memory caching
│   ├── database/          # Database repository
│   ├── models/            # Data models
│   ├── osrs/              # OSRS Wiki API client
│   └── worker/            # Background workers
├── migrations/            # Database migrations
├── cmd/
│   └── migrate/          # Migration tool
└── .github/
    └── workflows/        # CI/CD pipelines
```

## CI/CD

GitHub Actions automatically builds executables for:
- Linux (amd64, arm64)
- macOS (amd64, arm64 - Apple Silicon)
- Windows (amd64)

Builds are triggered on every push to `main` and available in [Releases](../../releases).

## Deployment

### Railway
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and deploy
railway login
railway init
railway up
```

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o osrs-price-api main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/osrs-price-api .
CMD ["./osrs-price-api"]
```

### DigitalOcean/AWS/GCP
1. Download the appropriate binary from releases
2. Set up PostgreSQL database
3. Configure environment variables
4. Run the binary with systemd or supervisor

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Support

- **Issues**: [GitHub Issues](../../issues)
- **Documentation**: See `/docs` folder
- **API Docs**: See `API_DOCUMENTATION.md`

## Acknowledgments

- OSRS Wiki for providing the Real-time Prices API
- RuneScape community for inspiration