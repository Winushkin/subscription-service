// Package main является входной точкой в приложение
package main

import (
	"context"
	"errors"
	"fmt"

	_ "subscription-service/docs"
	"subscription-service/internal/config"
	"subscription-service/internal/handlers"
	"subscription-service/internal/logger"
	"subscription-service/internal/repository"
	"subscription-service/internal/storage"
	"subscription-service/migrator"

	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const devMode = true

// @title Subscriptions API
// @version 1.0
// @description REST API for user subscriptions aggregation
// @host localhost
// @BasePath /api/v1
func main() {
	ctx, err := logger.NewLoggerContext(context.Background(), devMode)
	if err != nil {
		return
	}
	log, ok := logger.GetLoggerFromCtx(ctx)
	if !ok {
		return
	}

	cfg, err := config.NewAppConfig()
	if err != nil {
		log.Error(ctx, "failed to config app", zap.Error(err))
	}
	log.Info(ctx, "Starting application", zap.String("port", cfg.ServerPort))

	pool, err := storage.NewPool(ctx, cfg.PostgresConfig)
	if err != nil {
		log.Error(ctx, "Failed to connect to database", zap.Error(err))
		panic(err)
	}
	defer pool.Close()

	if err = migrator.Up(pool); err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}
	log.Debug(ctx, "migrations Upped")

	log.Info(ctx, "Connected to database")
	repo := repository.NewRepository(pool)

	// Init router
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(func(c *gin.Context) {
		ctxWithLogger := logger.NewContextWithLogger(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctxWithLogger)
		c.Next()
	})
	engine.Use(cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := handlers.New(repo)
	api := engine.Group("/api/v1")
	{
		api.POST("/subscriptions", h.CreateSubscription)
		api.GET("/subscriptions/:id", h.GetSubscription)
		api.GET("/subscriptions", h.ListSubscriptions)
		api.PUT("/subscriptions/:id", h.UpdateSubscription)
		api.DELETE("/subscriptions/:id", h.DeleteSubscription)

		api.POST("/reports/cost", h.CalculateCostReport)
	}

	router, err := graceful.New(
		engine,
		graceful.WithAddr(fmt.Sprintf(":%s", cfg.ServerPort)),
		graceful.WithShutdownTimeout(graceful.DefaultShutdownTimeout),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize router: %w", err))
	}
	defer router.Close()

	log.Info(ctx, "starting server")

	if err = router.RunWithContext(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		panic(fmt.Errorf("failed to run router: %w", err))
	}
}
