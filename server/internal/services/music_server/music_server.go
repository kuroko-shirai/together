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

func New(
	ctx context.Context,
	cfg *config.Config,
) (*MusicServer, error) {
	publisher, err := pubsub.NewPublisher(
		cfg.MusicServer.Address,
	)
	if err != nil {
		return nil, err
	}

	storage := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
	)
	if err := storage.Ping(
		ctx,
	).Err(); err != nil {
		return nil, err
	}

	return &MusicServer{
		publisher: publisher,
		storage:   storage,
	}, nil
}

func (this *MusicServer) Run(
	ctx context.Context,
) error {
	if err := this.publisher.Run(ctx); err != nil {
		return err
	}

	var eproc error
	for {
		status := this.storage.GetDel(
			ctx,
			utils.RedisKeyTrack,
		)
		if status.Err() != nil {
			if status.Err() != redis.Nil {
				eproc = status.Err()

				break
			}
			continue
		}

		cmd, err := status.Uint64()
		if err != nil {
			eproc = err
			continue
		}

		this.publisher.SendMessage(
			ctx,
			&pb.Message{
				Command: cmd,
			},
		)
	}

	return eproc
}

func (this *MusicServer) Stop(ctx context.Context) error {
	this.publisher.Stop(ctx)

	return nil
}
