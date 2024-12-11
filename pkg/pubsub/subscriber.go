package pubsub

import (
	"context"
	"io"

	"github.com/google/uuid"
	pb "github.com/kuroko-shirai/together/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Subscriber struct {
	connection *grpc.ClientConn
	client     pb.PublisherClient
	stream     grpc.ServerStreamingClient[pb.Message]
}

func NewSubscriber(
	ctx context.Context,
	addr string,
) (*Subscriber, error) {
	connection, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewPublisherClient(connection)
	stream, err := client.Subscribe(
		ctx,
		&pb.SubscribeRequest{
			ClientId: uuid.New().String(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		connection: connection,
		client:     client,
		stream:     stream,
	}, nil
}

func (this *Subscriber) Recv(
	handle func(msg *pb.Message) error,
) error {
	msg, err := this.stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}

		return err
	}

	return handle(msg)
}

func (this *Subscriber) Stop() error {
	if err := this.stream.CloseSend(); err != nil {
		return err
	}

	return this.connection.Close()
}
