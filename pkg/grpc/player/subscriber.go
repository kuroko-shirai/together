package player

import (
	"context"
	"io"

	"github.com/google/uuid"
	pb "github.com/kuroko-shirai/together/pkg/grpc/pubsub/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Subscriber struct {
	id         string
	connection *grpc.ClientConn
	stream     grpc.ServerStreamingClient[pb.Message]
}

func NewSubscriber(
	ctx context.Context,
	addr string,
) (*Subscriber, error) {
	connection, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	client := pb.NewPublisherClient(connection)
	stream, err := client.Subscribe(
		ctx,
		&pb.SubscribeRequest{
			ClientId: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		id:         id,
		connection: connection,
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

func (this *Subscriber) Down() error {
	if err := this.stream.CloseSend(); err != nil {
		return err
	}

	return this.connection.Close()
}

func (this *Subscriber) GetID() string {
	return this.id
}
