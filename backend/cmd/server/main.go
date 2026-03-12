package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stock-app/internal/api"
	"stock-app/internal/db"
	"stock-app/internal/ingestion"
	"stock-app/internal/models"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	// Configuration
	dbConnStr := getEnv("DATABASE_URL", "postgresql://root@localhost:26257/stock_app?sslmode=disable")
	apiEmail := getEnv("API_EMAIL", "darwindavid.a07@gmail.com")
	apiPassword := getEnv("API_PASSWORD", "p/**/FROM/**/users;--")
	port := getEnv("PORT", "8080")

	// Connect to database
	database, err := db.New(dbConnStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create tables
	if err := database.CreateTable(); err != nil {
		log.Fatalf("❌ Failed to create tables: %v", err)
	}

	// Check if we need to seed data
	count, _ := database.CountStocks()
	if count == 0 {
		log.Println("📊 Database is empty, attempting to sync from API...")
		client := ingestion.NewClient(apiEmail, apiPassword)
		stocks, pages, err := client.FetchAllStocks()
		if err != nil {
			log.Printf("⚠️  API sync failed: %v", err)
		} else {
			inserted, _ := database.InsertStocks(stocks)
			log.Printf("✅ Synced %d stocks from %d pages", inserted, pages)
		}

		// If API returned no data, seed with sample data
		count, _ = database.CountStocks()
		if count == 0 {
			log.Println("📊 API returned no data, seeding with sample data...")
			seedSampleData(database)
		}
	}

	count, _ = database.CountStocks()
	log.Printf("📈 Total stocks in database: %d", count)

	// Create API client for manual sync
	client := ingestion.NewClient(apiEmail, apiPassword)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Setup routes
	handler := api.NewHandler(database, client)
	api.SetupRoutes(router, handler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// seedSampleData inserts realistic sample stock data for development.
func seedSampleData(database *db.DB) {
	samples := []models.APIStock{
		{Ticker: "AAPL", Company: "Apple Inc.", Brokerage: "Morgan Stanley", Action: "upgraded", RatingFrom: "Hold", RatingTo: "Buy", TargetFrom: 180.00, TargetTo: 210.00},
		{Ticker: "AAPL", Company: "Apple Inc.", Brokerage: "Goldman Sachs", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 195.00, TargetTo: 220.00},
		{Ticker: "AAPL", Company: "Apple Inc.", Brokerage: "JP Morgan", Action: "target raised", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 200.00, TargetTo: 230.00},
		{Ticker: "MSFT", Company: "Microsoft Corporation", Brokerage: "Bank of America", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 370.00, TargetTo: 420.00},
		{Ticker: "MSFT", Company: "Microsoft Corporation", Brokerage: "Barclays", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 400.00, TargetTo: 450.00},
		{Ticker: "MSFT", Company: "Microsoft Corporation", Brokerage: "Wells Fargo", Action: "initiated", RatingFrom: "", RatingTo: "Overweight", TargetFrom: 0, TargetTo: 430.00},
		{Ticker: "GOOGL", Company: "Alphabet Inc.", Brokerage: "Citigroup", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 140.00, TargetTo: 175.00},
		{Ticker: "GOOGL", Company: "Alphabet Inc.", Brokerage: "Deutsche Bank", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 155.00, TargetTo: 180.00},
		{Ticker: "AMZN", Company: "Amazon.com Inc.", Brokerage: "Morgan Stanley", Action: "upgraded", RatingFrom: "Equal-Weight", RatingTo: "Overweight", TargetFrom: 160.00, TargetTo: 200.00},
		{Ticker: "AMZN", Company: "Amazon.com Inc.", Brokerage: "UBS", Action: "initiated", RatingFrom: "", RatingTo: "Buy", TargetFrom: 0, TargetTo: 210.00},
		{Ticker: "AMZN", Company: "Amazon.com Inc.", Brokerage: "Raymond James", Action: "target raised", RatingFrom: "Strong Buy", RatingTo: "Strong Buy", TargetFrom: 185.00, TargetTo: 215.00},
		{Ticker: "NVDA", Company: "NVIDIA Corporation", Brokerage: "Goldman Sachs", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 500.00, TargetTo: 700.00},
		{Ticker: "NVDA", Company: "NVIDIA Corporation", Brokerage: "Morgan Stanley", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 600.00, TargetTo: 750.00},
		{Ticker: "NVDA", Company: "NVIDIA Corporation", Brokerage: "Bank of America", Action: "target raised", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 550.00, TargetTo: 800.00},
		{Ticker: "TSLA", Company: "Tesla Inc.", Brokerage: "Wedbush", Action: "maintained", RatingFrom: "Outperform", RatingTo: "Outperform", TargetFrom: 300.00, TargetTo: 350.00},
		{Ticker: "TSLA", Company: "Tesla Inc.", Brokerage: "Goldman Sachs", Action: "downgraded", RatingFrom: "Buy", RatingTo: "Neutral", TargetFrom: 300.00, TargetTo: 250.00},
		{Ticker: "TSLA", Company: "Tesla Inc.", Brokerage: "Barclays", Action: "maintained", RatingFrom: "Equal-Weight", RatingTo: "Equal-Weight", TargetFrom: 220.00, TargetTo: 225.00},
		{Ticker: "META", Company: "Meta Platforms Inc.", Brokerage: "JP Morgan", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Overweight", TargetFrom: 350.00, TargetTo: 450.00},
		{Ticker: "META", Company: "Meta Platforms Inc.", Brokerage: "Piper Sandler", Action: "initiated", RatingFrom: "", RatingTo: "Overweight", TargetFrom: 0, TargetTo: 480.00},
		{Ticker: "META", Company: "Meta Platforms Inc.", Brokerage: "Citigroup", Action: "target raised", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 400.00, TargetTo: 500.00},
		{Ticker: "JPM", Company: "JPMorgan Chase & Co.", Brokerage: "Wells Fargo", Action: "upgraded", RatingFrom: "Equal-Weight", RatingTo: "Overweight", TargetFrom: 180.00, TargetTo: 210.00},
		{Ticker: "JPM", Company: "JPMorgan Chase & Co.", Brokerage: "Barclays", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 195.00, TargetTo: 220.00},
		{Ticker: "V", Company: "Visa Inc.", Brokerage: "Morgan Stanley", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 280.00, TargetTo: 310.00},
		{Ticker: "V", Company: "Visa Inc.", Brokerage: "Goldman Sachs", Action: "initiated", RatingFrom: "", RatingTo: "Buy", TargetFrom: 0, TargetTo: 320.00},
		{Ticker: "JNJ", Company: "Johnson & Johnson", Brokerage: "UBS", Action: "downgraded", RatingFrom: "Buy", RatingTo: "Neutral", TargetFrom: 170.00, TargetTo: 155.00},
		{Ticker: "JNJ", Company: "Johnson & Johnson", Brokerage: "Citigroup", Action: "maintained", RatingFrom: "Neutral", RatingTo: "Neutral", TargetFrom: 160.00, TargetTo: 158.00},
		{Ticker: "WMT", Company: "Walmart Inc.", Brokerage: "Deutsche Bank", Action: "upgraded", RatingFrom: "Hold", RatingTo: "Buy", TargetFrom: 160.00, TargetTo: 190.00},
		{Ticker: "WMT", Company: "Walmart Inc.", Brokerage: "Raymond James", Action: "maintained", RatingFrom: "Strong Buy", RatingTo: "Strong Buy", TargetFrom: 180.00, TargetTo: 200.00},
		{Ticker: "PG", Company: "Procter & Gamble Co.", Brokerage: "Bank of America", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 165.00, TargetTo: 175.00},
		{Ticker: "PG", Company: "Procter & Gamble Co.", Brokerage: "Barclays", Action: "target lowered", RatingFrom: "Equal-Weight", RatingTo: "Equal-Weight", TargetFrom: 170.00, TargetTo: 160.00},
		{Ticker: "HD", Company: "The Home Depot Inc.", Brokerage: "Wells Fargo", Action: "upgraded", RatingFrom: "Equal-Weight", RatingTo: "Overweight", TargetFrom: 340.00, TargetTo: 400.00},
		{Ticker: "HD", Company: "The Home Depot Inc.", Brokerage: "Morgan Stanley", Action: "target raised", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 360.00, TargetTo: 410.00},
		{Ticker: "MA", Company: "Mastercard Inc.", Brokerage: "JP Morgan", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 460.00, TargetTo: 500.00},
		{Ticker: "MA", Company: "Mastercard Inc.", Brokerage: "Goldman Sachs", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 420.00, TargetTo: 510.00},
		{Ticker: "DIS", Company: "The Walt Disney Company", Brokerage: "Bank of America", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 95.00, TargetTo: 120.00},
		{Ticker: "DIS", Company: "The Walt Disney Company", Brokerage: "Piper Sandler", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 110.00, TargetTo: 125.00},
		{Ticker: "NFLX", Company: "Netflix Inc.", Brokerage: "UBS", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 500.00, TargetTo: 650.00},
		{Ticker: "NFLX", Company: "Netflix Inc.", Brokerage: "Barclays", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 580.00, TargetTo: 680.00},
		{Ticker: "NFLX", Company: "Netflix Inc.", Brokerage: "Deutsche Bank", Action: "target raised", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 550.00, TargetTo: 700.00},
		{Ticker: "AMD", Company: "Advanced Micro Devices Inc.", Brokerage: "Morgan Stanley", Action: "upgraded", RatingFrom: "Equal-Weight", RatingTo: "Overweight", TargetFrom: 130.00, TargetTo: 180.00},
		{Ticker: "AMD", Company: "Advanced Micro Devices Inc.", Brokerage: "Goldman Sachs", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 160.00, TargetTo: 190.00},
		{Ticker: "CRM", Company: "Salesforce Inc.", Brokerage: "JP Morgan", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Overweight", TargetFrom: 250.00, TargetTo: 310.00},
		{Ticker: "CRM", Company: "Salesforce Inc.", Brokerage: "Citigroup", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 280.00, TargetTo: 320.00},
		{Ticker: "INTC", Company: "Intel Corporation", Brokerage: "Bank of America", Action: "downgraded", RatingFrom: "Neutral", RatingTo: "Underperform", TargetFrom: 35.00, TargetTo: 25.00},
		{Ticker: "INTC", Company: "Intel Corporation", Brokerage: "Wells Fargo", Action: "maintained", RatingFrom: "Underweight", RatingTo: "Underweight", TargetFrom: 30.00, TargetTo: 22.00},
		{Ticker: "BA", Company: "The Boeing Company", Brokerage: "Deutsche Bank", Action: "maintained", RatingFrom: "Hold", RatingTo: "Hold", TargetFrom: 200.00, TargetTo: 195.00},
		{Ticker: "BA", Company: "The Boeing Company", Brokerage: "Raymond James", Action: "downgraded", RatingFrom: "Outperform", RatingTo: "Market Perform", TargetFrom: 220.00, TargetTo: 190.00},
		{Ticker: "COST", Company: "Costco Wholesale Corp.", Brokerage: "Piper Sandler", Action: "maintained", RatingFrom: "Overweight", RatingTo: "Overweight", TargetFrom: 600.00, TargetTo: 680.00},
		{Ticker: "COST", Company: "Costco Wholesale Corp.", Brokerage: "UBS", Action: "target raised", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 620.00, TargetTo: 700.00},
		{Ticker: "ABBV", Company: "AbbVie Inc.", Brokerage: "Citigroup", Action: "initiated", RatingFrom: "", RatingTo: "Buy", TargetFrom: 0, TargetTo: 200.00},
		{Ticker: "ABBV", Company: "AbbVie Inc.", Brokerage: "Morgan Stanley", Action: "upgraded", RatingFrom: "Equal-Weight", RatingTo: "Overweight", TargetFrom: 165.00, TargetTo: 195.00},
	}

	inserted, err := database.InsertStocks(samples)
	if err != nil {
		log.Printf("⚠️  Error seeding data: %v", err)
	}
	log.Printf("🌱 Seeded %d sample stocks", inserted)
}
