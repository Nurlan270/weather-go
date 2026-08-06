package main

import (
	"errors"
	"github.com/Nurlan270/weather-go/internal/app"
	"log"
	"net/http"
)

func main() {
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("app finished unexpectedly: %v", err)
	}
}
