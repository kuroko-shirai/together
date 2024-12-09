package music_server

import (
	"context"
	"log"
	"net"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"google.golang.org/grpc"
)

type MusicServer struct {
	pb.UnimplementedPublisherServer

	server      *grpc.Server
	listener    net.Listener
	subscribers map[string]pb.Publisher_SubscribeServer
}

func New(config *Config) (*MusicServer, error) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	pb.RegisterPublisherServer(server, &MusicServer{
		subscribers: make(map[string]pb.Publisher_SubscribeServer),
	})

	return &MusicServer{
		server:   server,
		listener: listener,
	}, nil
}

func (s *MusicServer) Run() {
	s.server.Serve(s.listener)

}

func (s *MusicServer) Stop() {
	s.server.Stop()
}

func (s *MusicServer) SendMessage(ctx context.Context, req *pb.Message) (*pb.Response, error) {
	log.Printf("Received message: %s", req.Text)
	for _, subscriber := range s.subscribers {
		err := subscriber.SendMsg(&pb.Message{Text: req.Text})
		if err != nil {
			log.Printf("Error sending message to subscriber: %v", err)
		}
	}
	return &pb.Response{Result: "Message sent"}, nil
}

func (s *MusicServer) Subscribe(req *pb.SubscribeRequest, stream pb.Publisher_SubscribeServer) error {
	log.Printf("New subscriber: %s", req.ClientId)
	s.subscribers[req.ClientId] = stream
	return nil
}
