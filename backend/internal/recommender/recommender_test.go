package recommender

import (
	"testing"

	"stock-app/internal/models"
)

func TestScoreStock_Upgraded(t *testing.T) {
	stock := models.Stock{
		Ticker:     "AAPL",
		Company:    "Apple Inc.",
		Brokerage:  "Goldman Sachs",
		Action:     "upgraded",
		RatingFrom: "Hold",
		RatingTo:   "Buy",
		TargetFrom: 180.00,
		TargetTo:   210.00,
	}
	score, reason := ScoreStock(stock)
	if score <= 0 {
		t.Errorf("expected positive score for upgraded stock, got %f", score)
	}
	if reason == "" {
		t.Error("expected non-empty reason")
	}
	t.Logf("Upgraded stock score: %.2f, reason: %s", score, reason)
}

func TestScoreStock_Downgraded(t *testing.T) {
	stock := models.Stock{
		Ticker:     "INTC",
		Company:    "Intel Corporation",
		Action:     "downgraded",
		RatingFrom: "Buy",
		RatingTo:   "Sell",
		TargetFrom: 100.00,
		TargetTo:   60.00,
	}
	score, _ := ScoreStock(stock)

	upgradedStock := models.Stock{
		Ticker:     "AAPL",
		Company:    "Apple Inc.",
		Action:     "upgraded",
		RatingFrom: "Hold",
		RatingTo:   "Buy",
		TargetFrom: 180.00,
		TargetTo:   210.00,
	}
	upgradedScore, _ := ScoreStock(upgradedStock)

	if score >= upgradedScore {
		t.Errorf("downgraded stock score (%f) should be less than upgraded (%f)", score, upgradedScore)
	}
	t.Logf("Downgraded=%f vs Upgraded=%f", score, upgradedScore)
}

func TestScoreStock_Maintained(t *testing.T) {
	stock := models.Stock{
		Ticker:     "MSFT",
		Company:    "Microsoft",
		Action:     "maintained",
		RatingFrom: "Buy",
		RatingTo:   "Buy",
		TargetFrom: 400.00,
		TargetTo:   420.00,
	}
	score, _ := ScoreStock(stock)
	t.Logf("Maintained stock score: %.2f", score)
	// Maintained with same rating should be moderate score
	if score < 0 {
		t.Errorf("expected non-negative score for maintained Buy, got %f", score)
	}
}

func TestGetRecommendations_TopN(t *testing.T) {
	stocks := []models.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple", Action: "upgraded", RatingFrom: "Hold", RatingTo: "Strong Buy", TargetFrom: 180, TargetTo: 250},
		{ID: 2, Ticker: "INTC", Company: "Intel", Action: "downgraded", RatingFrom: "Buy", RatingTo: "Sell", TargetFrom: 100, TargetTo: 60},
		{ID: 3, Ticker: "MSFT", Company: "Microsoft", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 400, TargetTo: 420},
		{ID: 4, Ticker: "NVDA", Company: "NVIDIA", Action: "upgraded", RatingFrom: "Neutral", RatingTo: "Buy", TargetFrom: 500, TargetTo: 700},
		{ID: 5, Ticker: "TSLA", Company: "Tesla", Action: "downgraded", RatingFrom: "Buy", RatingTo: "Neutral", TargetFrom: 300, TargetTo: 250},
	}

	recs := GetRecommendations(stocks, 3)

	if len(recs) != 3 {
		t.Errorf("expected 3 recommendations, got %d", len(recs))
	}

	// Verify sorted by score descending
	for i := 1; i < len(recs); i++ {
		if recs[i].Score > recs[i-1].Score {
			t.Errorf("recommendations not sorted: index %d (%f) > index %d (%f)", i, recs[i].Score, i-1, recs[i-1].Score)
		}
	}

	// Top recommendation should be AAPL or NVDA (both upgraded)
	top := recs[0].Stock.Ticker
	if top != "AAPL" && top != "NVDA" {
		t.Errorf("expected AAPL or NVDA as top recommendation, got %s", top)
	}

	t.Logf("Top 3 recommendations:")
	for i, r := range recs {
		t.Logf("  %d. %s (score: %.2f): %s", i+1, r.Stock.Ticker, r.Score, r.Reason)
	}
}

func TestGetRecommendations_Empty(t *testing.T) {
	recs := GetRecommendations([]models.Stock{}, 5)
	if len(recs) != 0 {
		t.Errorf("expected 0 recommendations for empty input, got %d", len(recs))
	}
}

func TestGetRecommendations_AggregatesPerTicker(t *testing.T) {
	// Multiple entries for same ticker should be aggregated
	stocks := []models.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple", Brokerage: "Goldman", Action: "upgraded", RatingFrom: "Hold", RatingTo: "Buy", TargetFrom: 180, TargetTo: 210},
		{ID: 2, Ticker: "AAPL", Company: "Apple", Brokerage: "Morgan", Action: "maintained", RatingFrom: "Buy", RatingTo: "Buy", TargetFrom: 195, TargetTo: 220},
	}

	recs := GetRecommendations(stocks, 10)

	if len(recs) != 1 {
		t.Errorf("expected 1 recommendation (aggregated), got %d", len(recs))
	}

	if recs[0].Stock.Ticker != "AAPL" {
		t.Errorf("expected AAPL, got %s", recs[0].Stock.Ticker)
	}
}
