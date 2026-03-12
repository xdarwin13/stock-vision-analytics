package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"stock-app/internal/models"
)

// DB wraps the sql.DB connection.
type DB struct {
	conn *sql.DB
}

// New creates a new database connection.
func New(connStr string) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Connected to CockroachDB")
	return &DB{conn: conn}, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// CreateTable creates the stocks table if it doesn't exist.
func (d *DB) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS stocks (
			id SERIAL PRIMARY KEY,
			ticker STRING NOT NULL,
			company STRING NOT NULL,
			brokerage STRING NOT NULL,
			action STRING NOT NULL,
			rating_from STRING NOT NULL DEFAULT '',
			rating_to STRING NOT NULL DEFAULT '',
			target_from FLOAT NOT NULL DEFAULT 0,
			target_to FLOAT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			UNIQUE (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
		);
	`
	_, err := d.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	log.Println("Stocks table ready")
	return nil
}

// InsertStocks inserts or ignores duplicate stocks.
func (d *DB) InsertStocks(stocks []models.APIStock) (int, error) {
	if len(stocks) == 0 {
		return 0, nil
	}

	inserted := 0
	for _, s := range stocks {
		_, err := d.conn.Exec(`
			INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
			DO NOTHING
		`, s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo)
		if err != nil {
			log.Printf("Failed to insert stock %s: %v", s.Ticker, err)
			continue
		}
		inserted++
	}
	return inserted, nil
}

// CountStocks returns the total number of stocks.
func (d *DB) CountStocks() (int, error) {
	var count int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&count)
	return count, err
}

// GetStocks retrieves stocks with pagination, search, and sorting.
func (d *DB) GetStocks(search, sortBy, sortOrder string, page, pageSize int) (*models.StockListResponse, error) {
	// Build WHERE clause
	whereClause := ""
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		whereClause = fmt.Sprintf(` WHERE LOWER(ticker) LIKE $%d OR LOWER(company) LIKE $%d OR LOWER(brokerage) LIKE $%d OR LOWER(action) LIKE $%d`, argIdx, argIdx+1, argIdx+2, argIdx+3)
		args = append(args, searchLower, searchLower, searchLower, searchLower)
		argIdx += 4
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM stocks" + whereClause
	var total int
	err := d.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count stocks: %w", err)
	}

	// Validate sort columns
	validSortColumns := map[string]bool{
		"ticker": true, "company": true, "brokerage": true,
		"action": true, "target_from": true, "target_to": true,
		"rating_from": true, "rating_to": true, "created_at": true, "id": true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Calculate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	totalPages := (total + pageSize - 1) / pageSize

	// Build query
	query := fmt.Sprintf(
		"SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at FROM stocks%s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sortBy, sortOrder, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, s)
	}

	return &models.StockListResponse{
		Items:      stocks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetStockByID retrieves a single stock by its ID.
func (d *DB) GetStockByID(id int) (*models.Stock, error) {
	var s models.Stock
	err := d.conn.QueryRow(
		"SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at FROM stocks WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
		&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetAllStocks retrieves all stocks for recommendation analysis.
func (d *DB) GetAllStocks() ([]models.Stock, error) {
	rows, err := d.conn.Query(
		"SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at FROM stocks",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}
