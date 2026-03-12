package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"stock-app/internal/db"
	"stock-app/internal/ingestion"
	"stock-app/internal/recommender"
)

// Handler holds dependencies for route handlers.
type Handler struct {
	DB     *db.DB
	Client *ingestion.Client
}

// NewHandler creates a new Handler.
func NewHandler(database *db.DB, client *ingestion.Client) *Handler {
	return &Handler{
		DB:     database,
		Client: client,
	}
}

// GetStocks handles GET /api/stocks
func (h *Handler) GetStocks(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.DB.GetStocks(search, sortBy, sortOrder, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		log.Printf("❌ Error fetching stocks: %v", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetStockByID handles GET /api/stocks/:id
func (h *Handler) GetStockByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	stock, err := h.DB.GetStockByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetRecommendations handles GET /api/recommendations
func (h *Handler) GetRecommendations(c *gin.Context) {
	topNStr := c.DefaultQuery("top", "10")
	topN, _ := strconv.Atoi(topNStr)
	if topN < 1 {
		topN = 10
	}

	stocks, err := h.DB.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		log.Printf("❌ Error fetching stocks for recommendations: %v", err)
		return
	}

	recs := recommender.GetRecommendations(stocks, topN)
	c.JSON(http.StatusOK, gin.H{
		"recommendations": recs,
		"total":           len(recs),
	})
}

// SyncStocks handles POST /api/sync
func (h *Handler) SyncStocks(c *gin.Context) {
	stocks, pages, err := h.Client.FetchAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch stocks from API",
			"details": err.Error(),
		})
		log.Printf("❌ Error syncing stocks: %v", err)
		return
	}

	inserted, err := h.DB.InsertStocks(stocks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert stocks into database",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Synced %d stocks from %d API pages", inserted, pages)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Sync completed",
		"items_synced":  inserted,
		"items_fetched": len(stocks),
		"pages_scanned": pages,
	})
}

// GetStats handles GET /api/stats
func (h *Handler) GetStats(c *gin.Context) {
	count, err := h.DB.CountStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_stocks": count,
	})
}

// SetupRoutes configures all API routes.
func SetupRoutes(router *gin.Engine, handler *Handler) {
	api := router.Group("/api")
	{
		api.GET("/stocks", handler.GetStocks)
		api.GET("/stocks/:id", handler.GetStockByID)
		api.GET("/recommendations", handler.GetRecommendations)
		api.POST("/sync", handler.SyncStocks)
		api.GET("/stats", handler.GetStats)
	}
}
