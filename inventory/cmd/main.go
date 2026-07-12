package main

import (
	"context"
	"inventory/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
