package listener

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/kuroko-shirai/together/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Listener struct {
	server     *grpc.Server
	connection *grpc.ClientConn
	client     pb.PublisherClient
	stream     grpc.ServerStreamingClient[pb.Message]
}

func New(config *Config) (*Listener, error) {
	connection, err := grpc.NewClient(config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewPublisherClient(connection)
	stream, err := client.Subscribe(
		context.Background(),
		&pb.SubscribeRequest{
			ClientId: uuid.New().String(),
		},
	)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				continue
			}
			log.Printf("Received message: %s", resp)
		}
	}()

	time.Sleep(10 * time.Second)

	return &Listener{
		connection: connection,
		client:     client,
		stream:     stream,
		server:     grpc.NewServer(),
	}, nil
}

func (s *Listener) Run() {
	// for {
	// 	msg, err := s.stream.RecvMsg()
	// 	if err != nil {
	// 		continue
	// 	}
	// 	log.Printf("Received message: %s", msg.Text)
	// }
	// s.server.Serve(s.listener)
	// for {
	// 	s.stream.Recv()
	// }
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
