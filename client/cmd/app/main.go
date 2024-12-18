package main

import (
	"context"
	"log"

	"github.com/kuroko-shirai/together/client/internal/infra"
)

func main() {
	ctx := context.Background()

	app, err := infra.New(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer app.Down(ctx)

	app.Run(ctx)
}
