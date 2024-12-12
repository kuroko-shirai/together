package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
)

// Слушатель подписывается на рассылку сообщений от сервера,
// что позволяет множеству подписчиков получать актуальные
// синхронные обновления статуса сервера.
// С другой стороны, слушатель может послать команду
// серверу, чтобы обновить его статус.
type Listener struct {
	pb.PublisherServer

	publisher  *pubsub.Publisher
	subscriber *pubsub.Subscriber
}

func New(config *config.Config) (*Listener, error) {
	subscriber, err := pubsub.NewSubscriber(
		context.Background(),
		config.MusicServer.Address,
	)
	if err != nil {
		return nil, err
	}

	publisher, err := pubsub.NewPublisher(
		config.Listener.Address,
	)
	if err != nil {
		return nil, err
	}

	return &Listener{
		subscriber: subscriber,
		publisher:  publisher,
	}, nil
}

func (this *Listener) Run() error {
	var eproc error
	for {
		if err := this.subscriber.Recv(
			func(msg *pb.Message) error {
				fmt.Printf("client received ack: %s\n", msg)

				return nil
			},
		); err != nil {
			eproc = err
			break
		}
	}

	return eproc
}

func (this *Listener) Stop() error {
	return this.subscriber.Stop()
}

func (this *Listener) SendMessage(
	ctx context.Context,
	msg *pb.Message,
) (*pb.Response, error) {
	log.Printf("Received: %v", msg.GetText())
	return &pb.Response{Result: "sended: " + msg.GetText()}, nil
}
