package main

import (
	"log"

	"github.com/elvin-tajirzada/log-alerting/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Unable to create app: %v", err)
	}

	a.Start()
	a.Shutdown()
}
