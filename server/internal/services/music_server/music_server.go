package music_server

import (
	"context"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
)

// MusicServer - рассылает свой статус всем подписчикам.
// С другой стороны MusicServer принимает команды от
// подписчика, изменяет свой статуст и рассылает его всем.
type MusicServer struct {
	publisher  *pubsub.Publisher
	subscriber *pubsub.Subscriber
}

func New(config *config.Config) (*MusicServer, error) {
	publisher, err := pubsub.NewPublisher(config.MusicServer.Address)
	if err != nil {
		return nil, err
	}

	subscriber, err := pubsub.NewSubscriber(
		context.Background(),
		config.Listener.Address,
	)
	if err != nil {
		return nil, err
	}

	return &MusicServer{
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

func (this *MusicServer) Run() error {
	if err := this.publisher.Run(); err != nil {
		return err
	}

	var eproc error
	for {
		if err := this.subscriber.Recv(
			func(msg *pb.Message) error {
				this.publisher.SendMessage(
					context.TODO(),
					msg,
				)

				return nil
			},
		); err != nil {
			eproc = err
			break
		}
	}

	return eproc
}

func (this *MusicServer) Stop() error {
	this.publisher.Stop()

	return nil
}
