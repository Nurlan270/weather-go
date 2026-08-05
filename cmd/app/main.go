package main

import (
	"github.com/Nurlan270/weather-go/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app finished unexpectedly: %v", err)
	}
}
