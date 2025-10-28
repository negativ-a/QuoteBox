package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Adeel56/quotebox/internal/client"
	"github.com/Adeel56/quotebox/internal/db"
	"github.com/Adeel56/quotebox/internal/metrics"
	"github.com/Adeel56/quotebox/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QuoteHandler handles quote-related requests
type QuoteHandler struct {
	OpenRouterClient *client.OpenRouterClient
}

// NewQuoteHandler creates a new quote handler
func NewQuoteHandler(openRouterClient *client.OpenRouterClient) *QuoteHandler {
	return &QuoteHandler{
		OpenRouterClient: openRouterClient,
	}
}

// CreateQuoteRequest represents the request body for creating a quote
type CreateQuoteRequest struct {
	Tag       string `json:"tag" binding:"required"`
	Requestor string `json:"requestor"`
}

// QuoteResponse represents the response for a quote
type QuoteResponse struct {
	ID        uuid.UUID `json:"id"`
	Tag       string    `json:"tag"`
	Quote     string    `json:"quote"`
	Author    *string   `json:"author,omitempty"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateQuote handles POST /api/v1/quote
func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req CreateQuoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: fmt.Sprintf("Invalid JSON format: %v", err),
		})
		return
	}

	// Validate tag
	req.Tag = strings.TrimSpace(req.Tag)
	if req.Tag == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_tag",
			Message: "Tag cannot be empty",
		})
		return
	}

	if len(req.Tag) > 50 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_tag",
			Message: "Tag must be 50 characters or less",
		})
		return
	}

	// Record start time
	startTime := time.Now()

	// Generate quote from OpenRouter
	quoteText, err := h.OpenRouterClient.GenerateQuote(req.Tag)
	if err != nil {
		log.Printf("Error generating quote: %v", err)
		metrics.RecordQuoteError()
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "quote_generation_failed",
			Message: "Failed to generate quote. Please try again later.",
		})
		return
	}

	// Calculate latency
	latency := time.Since(startTime)
	latencyMs := int(latency.Milliseconds())

	// Record metrics
	metrics.RecordQuoteFetched(req.Tag)
	metrics.RecordLatency(latency.Seconds())

	// Determine tag source
	tagSource := models.GetTagSource(req.Tag)

	// Create quote record
	quote := models.Quote{
		Tag:       req.Tag,
		TagSource: tagSource,
		QuoteText: quoteText,
		Author:    nil, // OpenRouter doesn't typically return author
		Source:    "openrouter",
		CreatedAt: time.Now(),
		LatencyMs: latencyMs,
		ClientIP:  c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	// Save to database
	if err := db.DB.Create(&quote).Error; err != nil {
		log.Printf("Error saving quote to database: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "database_error",
			Message: "Failed to save quote",
		})
		return
	}

	log.Printf("Quote created successfully: ID=%s, Tag=%s, Latency=%dms", quote.ID, quote.Tag, quote.LatencyMs)

	// Return response
	c.JSON(http.StatusOK, QuoteResponse{
		ID:        quote.ID,
		Tag:       quote.Tag,
		Quote:     quote.QuoteText,
		Author:    quote.Author,
		Source:    quote.Source,
		CreatedAt: quote.CreatedAt,
	})
}

// GetQuotes handles GET /api/v1/quotes
func (h *QuoteHandler) GetQuotes(c *gin.Context) {
	tag := c.Query("tag")
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	query := db.DB.Order("created_at DESC").Limit(limit)

	if tag != "" {
		query = query.Where("tag = ?", tag)
	}

	var quotes []models.Quote
	if err := query.Find(&quotes).Error; err != nil {
		log.Printf("Error fetching quotes: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch quotes",
		})
		return
	}

	// Convert to response format
	responses := make([]QuoteResponse, len(quotes))
	for i, q := range quotes {
		responses[i] = QuoteResponse{
			ID:        q.ID,
			Tag:       q.Tag,
			Quote:     q.QuoteText,
			Author:    q.Author,
			Source:    q.Source,
			CreatedAt: q.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"quotes": responses,
		"count":  len(responses),
	})
}

// GetTags handles GET /api/v1/tags
func (h *QuoteHandler) GetTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tags": models.ValidTags,
	})
}
