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
	"net"
	"time"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"github.com/kuroko-shirai/together/pkg/pubsub"
	"google.golang.org/grpc"
)

type publisherServer struct {
	pb.UnimplementedPublisherServer

	server      *grpc.Server
	listener    net.Listener
	subscribers map[string]pb.Publisher_SubscribeServer
}

func (s *publisherServer) SendMessage(ctx context.Context, req *pb.Message) (*pb.Response, error) {
	log.Printf("Received message: %s", req.Text)
	for _, subscriber := range s.subscribers {
		err := subscriber.SendMsg(&pb.Message{Text: req.Text})
		if err != nil {
			log.Printf("Error sending message to subscriber: %v", err)
		}
	}
	return &pb.Response{Result: "Message sent"}, nil
}

func (s *publisherServer) Subscribe(req *pb.SubscribeRequest, stream pb.Publisher_SubscribeServer) error {
	log.Printf("New subscriber: %s", req.ClientId)
	s.subscribers[req.ClientId] = stream

	s.SendMessage(context.Background(), &pb.Message{
		Text: fmt.Sprintf("Message"),
	})

	return nil
}

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

	// lis, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// p := publisherServer{
	// 	subscribers: make(map[string]pb.Publisher_SubscribeServer),
	// }
	// srv := grpc.NewServer()
	// pb.RegisterPublisherServer(srv, &p)

	// fmt.Println("Server listening on port 8080")
	// if err := srv.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}
