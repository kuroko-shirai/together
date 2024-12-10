package main

import (
	"context"
	"log"

	"github.com/kuroko-shirai/together/client/internal/infra"
)

func main() {
	app, err := infra.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	app.Run(context.Background())
}
