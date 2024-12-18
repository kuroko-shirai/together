package music_server

import (
	"context"
	"log"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/grpc/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/grpc/pubsub/proto"
	"github.com/kuroko-shirai/together/utils"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
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

	log.Println("Server started on", cfg.MusicServer.Address)

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

	for {
		gem := this.storage.GetDel(
			ctx,
			utils.RedisKeyTrack,
		)
		if gem.Err() != nil {
			if gem.Err() != redis.Nil {
				log.Printf("error: %v", gem.Err())
			}

			continue
		}

		msg := &pb.Message{}
		if err := protojson.Unmarshal([]byte(gem.Val()), msg); err != nil {
			log.Printf("error: %v", err)

			continue
		}

		this.publisher.SendMessage(ctx, msg)
	}
}

func (this *MusicServer) Down(ctx context.Context) error {
	this.publisher.Down(ctx)

	return nil
}
