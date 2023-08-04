package main

import (
	"github.com/elvin-tacirzade/log-alerting/pkg/app"
	"log"
)

func main() {
	// create app
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to create a new app: %v", err)
	}

	// start app
	a.Start()

	// shutdown app
	a.Shutdown()
}
