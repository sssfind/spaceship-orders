package main

import (
	"context"
	"log"

	"payment/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize payment application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run payment application: %v", err)
	}
}
