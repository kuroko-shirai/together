package listener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kuroko-shirai/task"
	"github.com/kuroko-shirai/together/common/config"
	"github.com/kuroko-shirai/together/pkg/gramophone"
	"github.com/kuroko-shirai/together/pkg/grpc/player"
	pbplayer "github.com/kuroko-shirai/together/pkg/grpc/player/proto"
	pbpubsub "github.com/kuroko-shirai/together/pkg/grpc/pubsub/proto"
	"github.com/kuroko-shirai/together/utils"
	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

// Listener The listener subscribes to message broadcasts
// from the server, which allows multiple subscribers to
// receive up-to-date synchronous updates of the server
// status. On the other hand, the listener can send a
// command to the server to update its status.
type Listener struct {
	pbplayer.UnimplementedPlayerServer

	subscriber player.Subscriber
	storage    *redis.Client
	listener   *net.Listener

	gramophone gramophone.Gramophone
}

func New(
	ctx context.Context,
	cfg *config.Config,
) (*Listener, error) {
	subscriber, err := player.NewSubscriber(
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

		gramophone: *gramophone.New(),
	}, nil
}

func (this *Listener) Run(ctx context.Context) (err error) {
	this.start()

	for {
		if err = this.subscriber.Recv(
			func(msg *pbpubsub.Message) error {
				g := task.WithRecover(
					func(p any, args ...any) {
						log.Println("got a panic while processing command from music-server:", p)
					},
				)

				g.Do(
					func() func() error {
						return func() error {
							return this.process(msg)
						}
					}(),
				)

				return g.Wait()
			},
		); err != nil {
			log.Printf("got an error while processing command from music-server: %v", err)

			break
		}
	}

	return err
}

func (this *Listener) Down(context.Context) error {
	return this.subscriber.Down()
}

func (this *Listener) Play(
	ctx context.Context,
	msg *pbplayer.PlayRequest,
) (*pbplayer.PlayResponse, error) {
	status := this.storage.Set(
		ctx,
		utils.RedisKeyTrack, // stable key in the redis
		fmt.Sprintf(
			utils.MaskKeyRun,
			utils.CmdPlay,
			msg.GetAlbum(),
			msg.GetTitle(),
		), // send message from user
		10*time.Second, // record's ttl
	)
	if status.Err() != nil {
		return &pbplayer.PlayResponse{
			Result: utils.StatusError,
		}, status.Err()
	}

	return &pbplayer.PlayResponse{
		Result: utils.StatusOK,
	}, nil
}

func (this *Listener) Stop(
	ctx context.Context,
	msg *pbplayer.StopRequest,
) (*pbplayer.StopResponse, error) {
	status := this.storage.Set(
		ctx,
		utils.RedisKeyTrack, // stable key in the redis
		fmt.Sprintf(
			utils.MaskKeyStop,
			utils.CmdStop,
		), // send message from user
		10*time.Second, // record's ttl
	)
	if status.Err() != nil {
		fmt.Println(status.Err())
		return &pbplayer.StopResponse{
			Result: utils.StatusError,
		}, status.Err()
	}

	return &pbplayer.StopResponse{
		Result: utils.StatusOK,
	}, nil
}

func (this *Listener) GetListOfAlbums(
	ctx context.Context,
	msg *pbplayer.GetListOfAlbumsRequest,
) (*pbplayer.GetListOfAlbumsResponse, error) {
	files, err := os.ReadDir(utils.DirPlaylists)
	if err != nil {
		log.Fatal(err)
	}

	albums := make([]string, 0, len(files))
	for _, file := range files {
		albums = append(albums, file.Name())
	}

	return &pbplayer.GetListOfAlbumsResponse{
		Albums: albums,
	}, nil
}

func (this *Listener) GetListOfAlbumTracks(
	ctx context.Context,
	msg *pbplayer.GetListOfAlbumTracksRequest,
) (*pbplayer.GetListOfAlbumTracksResponse, error) {
	files, err := os.ReadDir(
		fmt.Sprintf(
			utils.DirPlaylistsAlbum,
			msg.GetAlbum(),
		),
	)
	if err != nil {
		return &pbplayer.GetListOfAlbumTracksResponse{
			Tracks: []string{},
		}, nil
	}

	tracks := make([]string, 0, len(files))
	for _, file := range files {
		tracks = append(tracks, file.Name())
	}

	return &pbplayer.GetListOfAlbumTracksResponse{
		Tracks: tracks,
	}, nil
}

// start registering and starting grpc server.
func (this *Listener) start() {
	go func() {
		server := grpc.NewServer()

		pbplayer.RegisterPlayerServer(server, this)

		if err := server.Serve(*this.listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
}

func (this *Listener) process(msg *pbpubsub.Message) error {
	switch msg.GetCommand() {
	case utils.CmdPlay:
		if msg.GetTrack().GetAlbum() == "" {
			return errors.New("invalid album")
		}

		if msg.GetTrack().GetAlbum() == "" {
			return errors.New("invalid track")
		}

		this.gramophone.Play(
			fmt.Sprintf(
				utils.DirPlaylistsAlbumTrack,
				msg.GetTrack().GetAlbum(),
				msg.GetTrack().GetTitle(),
			),
		)
	case utils.CmdStop:
		if this.gramophone.IsPlaying() {
			this.gramophone.Pause()
		}
	}

	return nil
}
