package app

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/Adeel56/quotebox/internal/app/handlers"
	"github.com/Adeel56/quotebox/internal/client"
	"github.com/Adeel56/quotebox/internal/db"
	"github.com/Adeel56/quotebox/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed frontend/*
var frontendFS embed.FS

// Server represents the application server
type Server struct {
	Router           *gin.Engine
	OpenRouterClient *client.OpenRouterClient
	QuoteHandler     *handlers.QuoteHandler
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Initialize metrics
	metrics.Init()

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize OpenRouter client
	openRouterClient := client.NewOpenRouterClient()

	// Create handlers
	quoteHandler := handlers.NewQuoteHandler(openRouterClient)

	// Create server
	server := &Server{
		OpenRouterClient: openRouterClient,
		QuoteHandler:     quoteHandler,
	}

	// Setup router
	server.setupRouter()

	return server
}

// setupRouter configures all routes
func (s *Server) setupRouter() {
	router := gin.Default()

	// Middleware for metrics
	router.Use(s.metricsMiddleware())

	// Health check
	router.GET("/healthz", s.healthCheck)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/quote", s.QuoteHandler.CreateQuote)
		apiV1.GET("/quotes", s.QuoteHandler.GetQuotes)
		apiV1.GET("/tags", s.QuoteHandler.GetTags)
	}

	// Serve frontend static files
	s.setupFrontend(router)

	s.Router = router
}

// setupFrontend configures frontend file serving
func (s *Server) setupFrontend(router *gin.Engine) {
	// Get frontend subdirectory
	frontendSubFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Printf("Warning: Could not load embedded frontend: %v", err)
		// Fallback to serving from filesystem
		if _, err := os.Stat("internal/app/frontend"); err == nil {
			router.Static("/", "internal/app/frontend")
		}
		return
	}

	// Serve embedded files
	router.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendSubFS, "index.html")
		if err != nil {
			log.Printf("Error reading index.html: %v", err)
			c.String(http.StatusInternalServerError, "Error loading page")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// Serve CSS
	router.GET("/style.css", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendSubFS, "style.css")
		if err != nil {
			c.String(http.StatusNotFound, "CSS not found")
			return
		}
		c.Data(http.StatusOK, "text/css; charset=utf-8", data)
	})

	// Serve JS
	router.GET("/app.js", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendSubFS, "app.js")
		if err != nil {
			c.String(http.StatusNotFound, "JS not found")
			return
		}
		c.Data(http.StatusOK, "application/javascript; charset=utf-8", data)
	})

	router.StaticFS("/static", http.FS(frontendSubFS))
}

// healthCheck handles GET /healthz
func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	if err := db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// metricsMiddleware records HTTP request metrics
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Record metrics after request is processed
		method := c.Request.Method
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		status := http.StatusText(c.Writer.Status())

		metrics.RecordHTTPRequest(method, route, status)
	}
}

// Run starts the server
func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	return s.Router.Run(":" + port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	return db.CloseDB()
}
