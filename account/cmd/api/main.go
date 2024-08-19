package main

import (
	"context"
	"log"

	"pragusga/internal/delivery/http"
	"pragusga/internal/events"
	repo "pragusga/internal/repo"
	"pragusga/internal/usecase"
	"pragusga/pkg/db"
	"pragusga/pkg/env"
	"pragusga/pkg/supertokens"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	cfg, err := env.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	// Initialize MongoDB
	mongoClient, err := db.NewMongoDBConnection(ctx, cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	database := mongoClient.Database(cfg.MongoDBDatabase)

	// Initialize Redis
	redisClient, err := db.NewRedisClient(ctx, cfg.RedisURI)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize SuperTokens
	err = supertokens.Init(supertokens.Config{
		ConnectionURI: cfg.SuperTokensConnectionURI,
		APIKey:        cfg.SuperTokensAPIKey,
		AppName:       cfg.SuperTokensAppName,
		APIDomain:     cfg.SuperTokensAPIDomain,
		WebsiteDomain: cfg.SuperTokensWebsiteDomain,
	})
	if err != nil {
		log.Fatalf("Failed to initialize SuperTokens: %v", err)
	}

	// Setup repositories
	userRepo := repo.NewUserRepository(database)

	// Setup event publisher
	userEventPublisher := events.NewUserEventPublisher(redisClient)

	// Setup use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, userEventPublisher)

	// Setup Gin router
	router := gin.Default()

	// Setup handlers
	http.NewAuthHandler(router, authUseCase, cfg)

	// Start server
	log.Fatal(router.Run(":" + cfg.PORT))
}
