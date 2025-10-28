package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Adeel56/quotebox/internal/metrics"
)

// OpenRouterClient handles API calls to OpenRouter
type OpenRouterClient struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient() *OpenRouterClient {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = "openrouter/auto"
	}

	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	return &OpenRouterClient{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatCompletionRequest represents the request to OpenRouter API
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response from OpenRouter API
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// APIError represents an error from the API
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// GenerateQuote generates a quote for the given tag
func (c *OpenRouterClient) GenerateQuote(tag string) (string, error) {
	prompt := fmt.Sprintf(
		"Generate a meaningful inspirational quote about %s. "+
			"The quote should be 1-2 sentences, insightful, and motivational. "+
			"Only return the quote text itself without any introduction or explanation.",
		tag,
	)

	request := ChatCompletionRequest{
		Model: c.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a wise philosopher who creates short, meaningful quotes. Always respond with only the quote text, nothing else.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.8,
		MaxTokens:   150,
	}

	// Try the request, with one retry on transient errors
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying OpenRouter API call (attempt %d)...", attempt+1)
			time.Sleep(500 * time.Millisecond)
		}

		quote, err := c.makeRequest(request)
		if err == nil {
			metrics.SetOpenRouterStatus(true)
			return quote, nil
		}

		lastErr = err

		// Check if error is retryable (429, 5xx)
		if !isRetryableError(err) {
			break
		}
	}

	metrics.SetOpenRouterStatus(false)
	return "", lastErr
}

// makeRequest makes the actual HTTP request to OpenRouter
func (c *OpenRouterClient) makeRequest(request ChatCompletionRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	log.Printf("Calling OpenRouter API: %s with model %s", url, c.Model)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenRouter API error: status=%d, body=%s", resp.StatusCode, string(body))
		return "", &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("API error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	quote := response.Choices[0].Message.Content
	
	// Validate quote is not empty or too short
	if len(quote) < 10 {
		log.Printf("Warning: Generated quote is too short or empty: '%s'", quote)
		return "", fmt.Errorf("generated quote is invalid or too short")
	}
	
	log.Printf("Successfully generated quote: %s", quote)

	return quote, nil
}

// HTTPError represents an HTTP error
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error) bool {
	if httpErr, ok := err.(*HTTPError); ok {
		// Retry on 429 (rate limit) and 5xx (server errors)
		return httpErr.StatusCode == 429 || (httpErr.StatusCode >= 500 && httpErr.StatusCode < 600)
	}
	return false
}
