# StockVision - Backend (Go)

This is the Go backend for the StockVision application. It provides a robust, concurrent API to manage, ingest, and analyze stock market ratings from external sources.

## Core Responsibilities
1. **API Server**: Serves a RESTful API using the `gin-gonic/gin` framework.
2. **Data Ingestion (`internal/ingestion`)**: Connects to the third-party KarenAI API to sync stock data.
3. **Database Layer (`internal/db`)**: Manages the connection and queries to CockroachDB using `lib/pq` (pure Go Postgres driver). Ensures data integrity with UPSERT operations.
4. **Recommendation Engine (`internal/recommender`)**: Evaluates stocks through a scoring algorithm based on analyst actions, rating momentum, and target prices to output top investment picks.

---

## Environment Variables Configuration

The backend relies on the `github.com/joho/godotenv` package to read configuration values seamlessly.

To configure your instance, create a `.env` file in this directory (`backend/`) from the provided example:
```bash
cp .env.example .env
```

### Required Variables:

| Variable | Description | Default Value |
| :--- | :--- | :--- |
| `PORT` | The port the Gin server will bind to. | `8080` |
| `DATABASE_URL` | The PostgreSQL/CockroachDB connection string. It must include user, host, port, dbname, and SSL mode. | `postgresql://root@localhost:26257/stock_app?sslmode=disable` |
| `API_TOKEN` | The JWT Authorization token used to bypass the login route of the external stock API. | `(Empty String)` |

*(Note: While `API_EMAIL` and `API_PASSWORD` were historically used, providing an updated `API_TOKEN` skips the authentication handshake entirely, increasing sync performance and avoiding repetitive login requests).*

---

## Technical Architecture

### 1. The Database (CockroachDB)
CockroachDB is chosen for distributed SQL capabilities. The schema relies on a single `stocks` table containing all relevant target prices, ratings, and analyst metadata securely indexed by `ticker`. When the backend scales, CockroachDB natively handles the distributed payload.

### 2. The Recommendation Algorithm
The recommendation strategy is fully contained in `/internal/recommender`. It iterates over the dataset and yields a score (-100 to +100) per stock based on:
- Positive/Negative sentiment of the 'action' (e.g. Upgrade vs Downgrade).
- Net difference in Target Prices.
- The weight of the Rating itself (Strong Buy > Equal-Weight > Underperform).

---

## Running the Application Locally

1. **Install Dependencies**
```bash
go mod tidy
```

2. **Run the Database**
Ensure CockroachDB is running and the database exists:
```bash
# SQL port (DB): 26257, HTTP UI port: 8081 (so it doesn't conflict with the Go API on 8080)
cockroach start-single-node \
  --insecure \
  --store=type=mem,size=0.25 \
  --listen-addr=localhost:26257 \
  --http-addr=localhost:8081 \
  --background

cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS stock_app;"
```

3. **Start the Server**
The server will automatically seed the database if it finds the table empty on startup.
```bash
go run ./cmd/server/main.go
```

## Testing
The application features comprehensive unit tests, especially focused on ensuring that the `recommender` algorithm is bulletproof.

```bash
go test ./internal/... -v
```
