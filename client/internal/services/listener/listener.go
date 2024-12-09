package listener

import (
	"context"
	"fmt"
	"log"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Listener struct {
	connection *grpc.ClientConn
	client     pb.PublisherClient
	stream     grpc.ServerStreamingClient[pb.Message]
}

func New(config *Config) (*Listener, error) {
	fmt.Println(">>>> 1")
	connection, err := grpc.NewClient(config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>> 2")

	client := pb.NewPublisherClient(connection)
	stream, err := client.Subscribe(
		context.Background(),
		&pb.SubscribeRequest{
			ClientId: "client1",
		},
	)
	if err != nil {
		return nil, err
	}

	return &Listener{
		connection: connection,
		client:     client,
		stream:     stream,
	}, nil
}

func (s *Listener) Run() {
	for {
		msg, err := s.stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		log.Printf("Received message: %s", msg.Text)
	}
}

func (s *Listener) Stop() {
	s.connection.Close()
}

func (s *Listener) Send() {
	_, err := s.client.SendMessage(context.Background(), &pb.Message{Text: "Hello, server!"})
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
}
