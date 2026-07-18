package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"iam/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("failed to run application: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	if err := application.Stop(); err != nil {
		log.Printf("error during graceful shutdown: %v", err)
	} else {
		log.Println("Application stopped gracefully")
	}
}
