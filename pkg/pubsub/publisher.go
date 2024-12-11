package pubsub

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

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

	msgChan chan *pb.Message
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
		msgChan:     make(chan *pb.Message),
	}

	pb.RegisterPublisherServer(server, &p)

	return &p, nil
}

func (this *Publisher) Run() error {
	go func() {
		if err := this.server.Serve(this.listener); err != nil {
		}
	}()

	return nil
}

func (this *Publisher) SendMessage(
	_ context.Context,
	msg *pb.Message,
) (*pb.Response, error) {
	this.msgChan <- msg

	return &pb.Response{
		Result: _ok,
	}, nil
}

func (this *Publisher) publish(msg *pb.Message) error {
	var eproc error

	for clientID, subscriber := range this.subscribers {
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

	return eproc
}

func (this *Publisher) Subscribe(
	req *pb.SubscribeRequest,
	stream pb.Publisher_SubscribeServer,
) error {
	log.Printf("Добавлен клиент %s", req.ClientId)
	this.subscribers[req.ClientId] = stream

	for {
		select {
		case msg := <-this.msgChan:
			if err := this.publish(msg); err != nil {
				log.Printf("Удаляется клиент %s", req.ClientId)
				delete(this.subscribers, req.ClientId)
			}
		}
	}
}

func (this *Publisher) Stop() {
	this.server.Stop()
}
