package main

import (
	"log"

	"github.com/mygallery/mygallery/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	application.Logger().Print(application.Banner())

	for _, warning := range application.Warnings() {
		application.Logger().Printf("⚠️ %s", warning)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
