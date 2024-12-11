package main

/*
import (
	"context"
	"log"

	"github.com/kuroko-shirai/together/server/internal/infra"
)

func main() {
	app, err := infra.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	app.Run(context.Background())
}
*/

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"github.com/kuroko-shirai/together/pkg/pubsub"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p, err := pubsub.NewPublisher(":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Stop()

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	for i := 0; ; i++ {
		p.SendMessage(ctx, &pb.Message{Text: fmt.Sprintf("message-%d", i)})
		time.Sleep(time.Second)
	}
}
