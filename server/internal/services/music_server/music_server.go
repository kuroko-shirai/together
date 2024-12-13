package music_server

import (
	"context"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
	"github.com/kuroko-shirai/together/utils"
	"github.com/redis/go-redis/v9"
)

// MusicServer: receives commands from a subscriber, and
// broadcasts commands to all subscribers, so that they,
// after processing the command, perform the corresponding
// actions with the track.
type MusicServer struct {
	publisher *pubsub.Publisher
	storage   *redis.Client
}

func New(config *config.Config) (*MusicServer, error) {
	publisher, err := pubsub.NewPublisher(
		config.MusicServer.Address,
	)
	if err != nil {
		return nil, err
	}

	storage := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err := storage.Ping(
		context.TODO(),
	).Err(); err != nil {
		return nil, err
	}

	return &MusicServer{
		publisher: publisher,
		storage:   storage,
	}, nil
}

func (this *MusicServer) Run() error {
	if err := this.publisher.Run(); err != nil {
		return err
	}

	var eproc error
	for {
		status := this.storage.GetDel(
			context.TODO(),
			utils.RedisKeyTrack,
		)
		if status.Err() != nil {
			if status.Err() != redis.Nil {
				eproc = status.Err()

				break
			}
			continue
		}

		this.publisher.SendMessage(
			context.TODO(),
			&pb.Message{Text: status.String()},
		)
	}

	return eproc
}

func (this *MusicServer) Stop() error {
	this.publisher.Stop()

	return nil
}
