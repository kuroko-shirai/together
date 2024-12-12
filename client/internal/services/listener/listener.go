package listener

import (
	"context"
	"fmt"

	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
)

type Listener struct {
	subscriber pubsub.Subscriber
}

func New(config *Config) (*Listener, error) {
	subscriber, err := pubsub.NewSubscriber(
		context.Background(),
		config.Address,
	)
	if err != nil {
		return nil, err
	}

	return &Listener{
		subscriber: *subscriber,
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
