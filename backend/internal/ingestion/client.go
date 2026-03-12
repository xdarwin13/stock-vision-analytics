package ingestion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"stock-app/internal/models"
)

const (
	baseURL  = "https://api.karenai.click/swechallenge"
	loginURL = baseURL + "/login"
	listURL  = baseURL + "/list"
)

// Client handles communication with the KarenAI stock API.
type Client struct {
	httpClient *http.Client
	email      string
	password   string
	token      string
}

// NewClient creates a new API client.
func NewClient(email, password string) *Client {
	return &Client{
		httpClient: &http.Client{},
		email:      email,
		password:   password,
	}
}

// Login authenticates with the API and retrieves a bearer token.
func (c *Client) Login() error {
	payload := models.LoginRequest{
		Username: c.email,
		Password: c.password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	resp, err := c.httpClient.Post(loginURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	var loginResp models.LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if loginResp.AuthToken == "" {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	c.token = loginResp.AuthToken
	log.Println("✅ Successfully logged in to KarenAI API")
	return nil
}

// FetchStocks fetches a single page of stocks from the API.
func (c *Client) FetchStocks(nextPage string) (*models.APIResponse, error) {
	reqURL := listURL
	if nextPage != "" {
		reqURL = fmt.Sprintf("%s?next_page=%s", listURL, url.QueryEscape(nextPage))
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &apiResp, nil
}

// FetchAllStocks fetches all pages of stocks from the API.
func (c *Client) FetchAllStocks() ([]models.APIStock, int, error) {
	if c.token == "" {
		if err := c.Login(); err != nil {
			return nil, 0, err
		}
	}

	var allStocks []models.APIStock
	nextPage := ""
	pages := 0

	for {
		resp, err := c.FetchStocks(nextPage)
		if err != nil {
			return nil, pages, fmt.Errorf("failed to fetch page %d: %w", pages+1, err)
		}
		pages++
		allStocks = append(allStocks, resp.Items...)

		log.Printf("📄 Page %d: fetched %d stocks (total: %d)", pages, len(resp.Items), len(allStocks))

		if resp.NextPage == "" {
			break
		}
		nextPage = resp.NextPage
	}

	return allStocks, pages, nil
}
