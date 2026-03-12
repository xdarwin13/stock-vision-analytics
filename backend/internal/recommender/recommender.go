package recommender

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"stock-app/internal/models"
)

// actionScores maps analyst actions to scores.
var actionScores = map[string]float64{
	"upgraded":    1.0,
	"initiated":   0.7,
	"reiterated":  0.5,
	"maintained":  0.3,
	"target raised": 0.6,
	"target lowered": -0.3,
	"downgraded":  -0.8,
}

// ratingScores maps analyst ratings to numeric values.
var ratingScores = map[string]float64{
	"strong buy":  5.0,
	"buy":         4.0,
	"outperform":  3.5,
	"overweight":  3.0,
	"market perform": 2.5,
	"hold":        2.0,
	"neutral":     2.0,
	"equal-weight": 2.0,
	"sector perform": 2.0,
	"underweight": 1.5,
	"underperform": 1.0,
	"sell":        0.5,
	"strong sell": 0.0,
}

// getRatingScore returns a numeric score for a rating string.
func getRatingScore(rating string) float64 {
	normalized := strings.ToLower(strings.TrimSpace(rating))
	if score, ok := ratingScores[normalized]; ok {
		return score
	}
	return 2.5 // default neutral score
}

// getActionScore returns a numeric score for an action string.
func getActionScore(action string) float64 {
	normalized := strings.ToLower(strings.TrimSpace(action))
	if score, ok := actionScores[normalized]; ok {
		return score
	}
	return 0.0
}

// ScoreStock calculates a composite recommendation score for a stock.
func ScoreStock(stock models.Stock) (float64, string) {
	score := 0.0
	reasons := []string{}

	// 1. Action score (weight: 30%)
	actionScore := getActionScore(stock.Action)
	score += actionScore * 30
	if actionScore > 0.5 {
		reasons = append(reasons, "Positive analyst action: "+stock.Action)
	} else if actionScore < 0 {
		reasons = append(reasons, "Negative analyst action: "+stock.Action)
	}

	// 2. Rating improvement (weight: 30%)
	ratingFrom := getRatingScore(stock.RatingFrom)
	ratingTo := getRatingScore(stock.RatingTo)
	ratingDelta := ratingTo - ratingFrom
	score += ratingDelta * 6 // Scale to ~30 points max
	if ratingDelta > 0 {
		reasons = append(reasons, "Rating improved: "+stock.RatingFrom+" → "+stock.RatingTo)
	} else if ratingTo >= 4.0 {
		reasons = append(reasons, "Strong rating: "+stock.RatingTo)
	}

	// 3. Target price change (weight: 25%)
	if stock.TargetFrom > 0 {
		targetChange := (stock.TargetTo - stock.TargetFrom) / stock.TargetFrom * 100
		// Cap to avoid outliers
		targetChange = math.Max(-50, math.Min(50, targetChange))
		score += targetChange * 0.5
		if targetChange > 5 {
			reasons = append(reasons, fmt.Sprintf("Target price increase: +%.1f%%", targetChange))
		}
	}

	// 4. Current rating strength (weight: 15%)
	currentRating := getRatingScore(stock.RatingTo)
	score += currentRating * 3
	if currentRating >= 4.0 {
		reasons = append(reasons, "Analyst rates as "+stock.RatingTo)
	}

	reason := strings.Join(reasons, "; ")
	if reason == "" {
		reason = "Neutral outlook based on available data"
	}

	return score, reason
}

// GetRecommendations returns the top N stock recommendations.
func GetRecommendations(stocks []models.Stock, topN int) []models.Recommendation {
	// Aggregate scores per ticker (average across brokerages)
	type tickerData struct {
		stocks []models.Stock
		scores []float64
	}
	tickerMap := map[string]*tickerData{}

	for _, s := range stocks {
		if _, ok := tickerMap[s.Ticker]; !ok {
			tickerMap[s.Ticker] = &tickerData{}
		}
		score, _ := ScoreStock(s)
		tickerMap[s.Ticker].stocks = append(tickerMap[s.Ticker].stocks, s)
		tickerMap[s.Ticker].scores = append(tickerMap[s.Ticker].scores, score)
	}

	// Create recommendations with averaged scores
	var recs []models.Recommendation
	for _, data := range tickerMap {
		avgScore := 0.0
		for _, s := range data.scores {
			avgScore += s
		}
		avgScore /= float64(len(data.scores))

		// Pick the best individual record for display
		bestIdx := 0
		bestScore := -math.MaxFloat64
		for i, s := range data.scores {
			if s > bestScore {
				bestScore = s
				bestIdx = i
			}
		}
		_, reason := ScoreStock(data.stocks[bestIdx])

		recs = append(recs, models.Recommendation{
			Stock:  data.stocks[bestIdx],
			Score:  math.Round(avgScore*100) / 100,
			Reason: reason,
		})
	}

	// Sort by score descending
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})

	if topN > 0 && topN < len(recs) {
		recs = recs[:topN]
	}

	return recs
}
