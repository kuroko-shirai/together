package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/pubsub"
	pb "github.com/kuroko-shirai/together/pkg/pubsub/proto"
	"github.com/kuroko-shirai/together/utils"
	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

// Listener: The listener subscribes to message broadcasts
// from the server, which allows multiple subscribers to
// receive up-to-date synchronous updates of the server
// status. On the other hand, the listener can send a
// command to the server to update its status.
type Listener struct {
	pb.UnimplementedPublisherServer

	subscriber pubsub.Subscriber
	storage    *redis.Client
	listener   *net.Listener
}

func New(
	ctx context.Context,
	cfg *config.Config,
) (*Listener, error) {
	subscriber, err := pubsub.NewSubscriber(
		context.Background(),
		cfg.MusicServer.Address,
	)
	if err != nil {
		return nil, err
	}

	storage := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := storage.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	listener, err := cfg.GetAvailableListener()
	if err != nil {
		return nil, err
	}

	log.Println("Client started on", listener.Addr())
	return &Listener{
		subscriber: *subscriber,
		storage:    storage,
		listener:   &listener,
	}, nil
}

func (this *Listener) Run(context.Context) error {
	var eproc error

	go func() {
		server := grpc.NewServer()

		pb.RegisterPublisherServer(server, this)

		if err := server.Serve(*this.listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	for {
		if err := this.subscriber.Recv(
			func(msg *pb.Message) error {
				fmt.Printf("client received ack: %s\n", msg)
				// TODO: Here we need to process the
				// received command and perform one of the
				// actions with the music track:
				// play/pause/next/prev

				return nil
			},
		); err != nil {
			eproc = err
			break
		}
	}

	return eproc
}

func (this *Listener) Stop(context.Context) error {
	return this.subscriber.Stop()
}

type Message struct {
	ID      string
	Message string
}

func (this *Listener) SendMessage(
	ctx context.Context,
	msg *pb.Message,
) (*pb.Response, error) {
	status := this.storage.Set(ctx, utils.RedisKeyTrack, msg.GetText(), 10*time.Second)
	if status.Err() != nil {
		return &pb.Response{
			Result: "error",
		}, status.Err()
	}

	return &pb.Response{
		Result: "ok",
	}, nil
}
