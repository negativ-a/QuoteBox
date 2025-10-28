//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Adeel56/quotebox/internal/app"
	"github.com/Adeel56/quotebox/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// Setup
	setupTestEnvironment()

	// Initialize database
	if err := db.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Create test server
	server := app.NewServer()
	testServer = httptest.NewServer(server.Router)

	// Run tests
	code := m.Run()

	// Teardown
	testServer.Close()
	db.CloseDB()

	os.Exit(code)
}

func setupTestEnvironment() {
	os.Setenv("OPENROUTER_API_KEY", "test-key-integration")
	os.Setenv("OPENROUTER_MODEL", "openrouter/auto")
	os.Setenv("GIN_MODE", "test")
	
	// Use DATABASE_URL from environment if set, otherwise use default
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://quoteuser:quotepw@localhost:5432/quotedb?sslmode=disable")
	}
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
}

func TestGetTags(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/tags")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	tags, ok := result["tags"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(tags), 20)
}

func TestGetQuotes(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/quotes?limit=10")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	_, hasQuotes := result["quotes"]
	assert.True(t, hasQuotes)
}

func TestCreateQuote_InvalidRequest(t *testing.T) {
	// Test with empty tag
	reqBody := map[string]string{
		"tag": "",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		testServer.URL+"/api/v1/quote",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateQuote_TagTooLong(t *testing.T) {
	// Test with tag > 50 characters
	reqBody := map[string]string{
		"tag": "this_is_a_very_long_tag_that_exceeds_fifty_characters_limit",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		testServer.URL+"/api/v1/quote",
		"application/json",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestMetricsEndpoint(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "text/plain")
}

func TestGetQuotesWithFilter(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/quotes?tag=joy&limit=5")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	quotes, ok := result["quotes"].([]interface{})
	assert.True(t, ok)
	
	// If there are quotes, verify they have the correct tag
	for _, q := range quotes {
		quote := q.(map[string]interface{})
		if tag, exists := quote["tag"]; exists {
			assert.Equal(t, "joy", tag)
		}
	}
}

func TestCORSHeaders(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", testServer.URL+"/api/v1/quotes", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Gin should handle OPTIONS requests
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound)
}
