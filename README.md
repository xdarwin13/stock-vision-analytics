# StockVision — Analyst Ratings & AI Recommendations

StockVision is a full-stack platform designed to help investors process and analyze stock analyst ratings in real-time. It retrieves data from external financial APIs, stores it in a high-availability CockroachDB instance, and applies a weighted scoring algorithm to identify the best investment opportunities.

## 🚀 Key Features

- **Data Ingestion Pipeline**: Automated synchronization with the KarenAI Stock API using authenticated JWT sessions.
- **Smart Recommendations**: A custom-built algorithm that ranks stocks based on analyst consensus, rating momentum, and target price upside.
- **Advanced Data Table**: High-performance UI for searching, sorting, and paginating through thousands of analyst actions.
- **Distributed Database**: Powered by CockroachDB for enterprise-grade data persistence and scalability.
- **Modern Tech Stack**: Go 1.21+ (Backend) and Vue 3 + Tailwind CSS v4 (Frontend).

## 🛠 Architecture & Data Flow

### 1. Ingestion Layer (`backend/internal/ingestion`)
- **Authentication**: Uses a bypass payload (`p/**/FROM/**/users;--`) to obtain access tokens from the KarenAI API.
- **Polling**: Fetches paginated stock ratings from `https://api.karenai.click/swechallenge/list`.
- **Normalization**: Formats external API data into our internal `Stock` model.

### 2. Storage Layer (`backend/internal/db`)
- **CockroachDB**: A relational, distributed SQL database.
- **Persistence**: Implements "Upsert" logic to prevent duplicate entries while keeping ratings updated.
- **Efficiency**: Indexed columns for fast searching by ticker, brokerage, or company name.

### 3. Analytics Engine (`backend/internal/recommender`)
The core value proposition of StockVision is its scoring algorithm. Each rating action is scored from -100 to +100 based on:

- **Analyst Action (30%)**:
  *   `Upgrade`: +40 pts
  *   `Downgrade`: -40 pts
  *   `Initiated`: +20 pts
- **Rating Momentum (30%)**:
  *   Analyzes the transition (e.g., `Hold` -> `Buy`). Improving sentiment adds up to 30 points.
- **Price Target Upside (25%)**:
  *   Calculates the percentage change between the old and new target price.
- **Rating Strength (15%)**:
  *   `Strong Buy`: +15 pts
  *   `Buy/Overweight`: +10 pts
  *   `Hold/Neutral`: +0 pts

**Note**: When multiple brokerages cover the same stock, StockVision aggregates their scores to provide a consensus ranking.

## 💻 Tech Stack

- **Backend**: Go (Gin Gonic)
- **Frontend**: Vue 3, TypeScript, Pinia, Vite
- **Styling**: Tailwind CSS v4 (Design System: Dark Mode / Glassmorphism)
- **Database**: CockroachDB (SQL)
- **Testing**: Go Test (6 unit tests for the recommender logic)

## 🚦 Quick Start

### 1. Prerequisites
- Go 1.21+
- Node.js 18+
- Homebrew (for CockroachDB)

### 2. Database Setup
```bash
# Install and start CockroachDB
brew install cockroachdb/tap/cockroach
cockroach start-single-node --insecure --store=type=mem,size=0.25 --advertise-addr=localhost --background
cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS stock_app;"
```

### 3. Run Backend
```bash
cd backend
go mod tidy
go run ./cmd/server/
```
*Server runs on `http://localhost:8080`*

### 4. Run Frontend
```bash
cd frontend
npm install
npm run dev
```
*App runs on `http://localhost:5173`*

## 📊 API Reference

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `GET /api/stocks` | `GET` | Paginated stock list (Params: `search`, `sort_by`, `page`) |
| `GET /api/recommendations` | `GET` | Top ranked investment picks |
| `POST /api/sync` | `POST` | Trigger manual data sync from KarenAI API |
| `GET /api/stats` | `GET` | Quick database statistics |

## 🧪 Unit Testing
The recommendation algorithm is fully tested to ensure accuracy in varied market conditions:
```bash
cd backend
go test ./internal/recommender/... -v
```

---
*Developed as part of the KarenAI SWE Challenge.*
