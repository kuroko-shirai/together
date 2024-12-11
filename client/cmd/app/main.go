package main

/*
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
*/

import (
	"context"
	"fmt"
	"log"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"github.com/kuroko-shirai/together/pkg/pubsub"
)

func main() {
	s, err := pubsub.NewSubscriber(context.Background(), ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer s.Stop()

	for {
		if err := s.Recv(func(msg *pb.Message) error {
			fmt.Printf("client received ack: %s\n", msg)

			return nil
		}); err != nil {
			fmt.Println(err)
			break
		}
	}
}
