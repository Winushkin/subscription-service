package main

import (
	"context"
	"fmt"
	"subscription-service/internal/config"
	"subscription-service/internal/storage"
	"subscription-service/migrator"
)

func main(){
	ctx := context.Background()
	cfg, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}
	pool, err := storage.NewPool(ctx, cfg.PostgresConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create postgres pool: %w", err))
	}

	if err = migrator.Down(pool); err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}
}