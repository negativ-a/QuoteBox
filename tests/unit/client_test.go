package unit

import (
	"testing"

	"github.com/Adeel56/quotebox/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestHTTPError_Error(t *testing.T) {
	err := &client.HTTPError{
		StatusCode: 500,
		Message:    "Internal Server Error",
	}

	expected := "HTTP 500: Internal Server Error"
	assert.Equal(t, expected, err.Error())
}

func TestNewOpenRouterClient(t *testing.T) {
	// Set required environment variable
	t.Setenv("OPENROUTER_API_KEY", "test-key")
	t.Setenv("OPENROUTER_MODEL", "test-model")
	t.Setenv("OPENROUTER_BASE_URL", "https://test.example.com")

	c := client.NewOpenRouterClient()

	assert.NotNil(t, c)
	assert.Equal(t, "test-key", c.APIKey)
	assert.Equal(t, "test-model", c.Model)
	assert.Equal(t, "https://test.example.com", c.BaseURL)
	assert.NotNil(t, c.HTTPClient)
}
