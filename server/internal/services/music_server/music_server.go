package music_server

import (
	"context"
	"fmt"
	"time"

	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
)

type MusicServer struct {
	publisher pubsub.Publisher
}

func New(config *Config) (*MusicServer, error) {
	publisher, err := pubsub.NewPublisher(config.Address)
	if err != nil {
		return nil, err
	}

	return &MusicServer{
		publisher: *publisher,
	}, nil
}

func (this *MusicServer) Run() error {
	if err := this.publisher.Run(); err != nil {
		return err
	}

	ctx := context.Background()
	for i := 0; i <= 60; i++ {
		this.publisher.SendMessage(ctx, &pb.Message{Text: fmt.Sprintf("message-%d", i)})
		time.Sleep(time.Second)
	}

	return nil
}

func (this *MusicServer) Stop() error {
	this.publisher.Stop()

	return nil
}
