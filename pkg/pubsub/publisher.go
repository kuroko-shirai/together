package pubsub

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	"google.golang.org/grpc"
)

const (
	_ok = "ok"

	_feclient = "[%s] %s"
	_ferrors  = "error: %s"

	_tcp = "tcp"
)

type Publisher struct {
	pb.UnimplementedPublisherServer

	server      *grpc.Server
	listener    net.Listener
	subscribers map[string]pb.Publisher_SubscribeServer
}

func NewPublisher(
	address string,
) (*Publisher, error) {
	listener, err := net.Listen(_tcp, address)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()

	p := Publisher{
		subscribers: make(map[string]pb.Publisher_SubscribeServer),
		listener:    listener,
		server:      server,
	}

	pb.RegisterPublisherServer(server, &p)

	return &p, nil
}

func (this *Publisher) Run() error {
	go func() {
		if err := this.server.Serve(this.listener); err != nil {
			//	return err
		}
	}()

	return nil
}

func (this *Publisher) SendMessage(
	ctx context.Context,
	msg *pb.Message,
) (*pb.Response, error) {
	var eproc error

	for clientID, subscriber := range this.subscribers {
		subscriber.Context()
		if err := subscriber.SendMsg(msg); err != nil {
			eclient := errors.New(
				fmt.Sprintf(_feclient,
					clientID,
					err.Error(),
				),
			)
			eproc = errors.Join(eproc, eclient)
		}
	}

	if eproc != nil {
		return &pb.Response{
			Result: fmt.Sprintf(_ferrors, eproc.Error()),
		}, eproc
	}

	return &pb.Response{
		Result: _ok,
	}, nil
}

func (this *Publisher) Subscribe(
	req *pb.SubscribeRequest,
	stream pb.Publisher_SubscribeServer,
) error {
	this.subscribers[req.ClientId] = stream

	for {
		time.Sleep(time.Second)
	}
}

func (this *Publisher) Stop() {
	this.server.Stop()
}
