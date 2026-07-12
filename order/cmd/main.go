package main

import (
	"context"
	"log"
	"order/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize order application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run order application: %v", err)
	}
}
