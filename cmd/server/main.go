// main.go
package main

import (
	"chat-websocket/api"
	"chat-websocket/config"
	"chat-websocket/db"
	"chat-websocket/redis"
	"chat-websocket/repository"
	"chat-websocket/service"
	"chat-websocket/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Load configuration.
	cfg := config.LoadConfig()

	// 2. Initialize MySQL database.
	dbConn := db.InitMySQL(cfg)
	// dbConn can be used by repositories.

	// 3. Initialize Redis client.
	redisClient := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	defer redisClient.Close()

	// 4. Initialize Redis Pub/Sub repository.
	pubSubRepo := redis.NewPubSubRepository(redisClient)

	// 5. Initialize repositories.
	messageRepo := repository.NewMessageRepository(dbConn)
	// For clients, here we use the MySQL-based repository (you can replace with your own implementation)
	_ = repository.NewClientRepository(dbConn)

	// 6. Initialize services.
	messageService := service.NewMessageService(pubSubRepo)
	_ = service.NewRoomService(pubSubRepo)

	// 7. Initialize use cases.
	roomUseCase := usecase.NewRoomUseCase(pubSubRepo)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, messageService)

	// 8. Initialize API router (pass both roomUseCase and messageUseCase).
	router := api.NewRouter(roomUseCase, messageUseCase)

	// 9. Start HTTP server.
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 10. Graceful shutdown.
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
