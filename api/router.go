package api

import (
	"chat-websocket/usecase"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter sets up the HTTP routes for the WebSocket chat service.
func NewRouter(roomUseCase *usecase.RoomUseCase, messageUseCase *usecase.MessageUseCase) *gin.Engine {
	router := gin.Default()

	// Create a new WebSocketHandler with the provided use cases.
	wsHandler := NewWebSocketHandler(roomUseCase, messageUseCase)

	// Define the route for WebSocket connections.
	router.GET("/chat", func(c *gin.Context) {
		wsHandler.HandleConnection(c.Writer, c.Request)
	})

	// Set up Prometheus metrics endpoint.
	router.GET("/metrics", prometheusHandler())

	return router
}

// prometheusHandler wraps promhttp.Handler() to make it compatible with Gin.
func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
