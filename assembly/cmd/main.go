package main

import (
	"context"
	"log"

	"assembly/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create assembly app: %v", err)
	}

	err = application.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run assembly app: %v", err)
	}
}
