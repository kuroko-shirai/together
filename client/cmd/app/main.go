package main

import (
	"context"
	"fmt"

	"github.com/kuroko-shirai/together/client/internal/infra"
)

func main() {
	app, err := infra.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	app.Run(context.TODO())
}
