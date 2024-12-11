package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/fnv"
	"net"

	pb "github.com/kuroko-shirai/together/pkg/proto"
	cmap "github.com/orcaman/concurrent-map/v2"
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
	subscribers cmap.ConcurrentMap[string, pb.Publisher_SubscribeServer]

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
		subscribers: cmap.NewWithCustomShardingFunction[string, pb.Publisher_SubscribeServer](func(key string) uint32 {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(key); err != nil {
				panic(err)
			}
			h := fnv.New32()
			h.Write(buf.Bytes())
			return h.Sum32()
		}),
		listener: listener,
		server:   server,
		msgChan:  make(chan *pb.Message),
	}

	return &p, nil
}

func (this *Publisher) Run() error {
	go func() {
		pb.RegisterPublisherServer(this.server, this)

		this.server.Serve(this.listener)
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

	for clientID, subscriber := range this.subscribers.Items() {
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
	this.subscribers.Set(req.ClientId, stream)

	for {
		select {
		case msg := <-this.msgChan:
			if err := this.publish(msg); err != nil {
				this.subscribers.Remove(req.ClientId)
			}
		}
	}
}

func (this *Publisher) Stop() {
	this.server.Stop()
}
