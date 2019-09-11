package bulbasaur

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"

	pb "github.com/lj-211/bulbasaur/protocol"
)

type PartnerInfo struct {
	PartnerId  uint64
	Addr       string
	LastActive time.Time
	Role       pb.Role
}

var Info *PartnerInfo = &PartnerInfo{
	PartnerId: uint64(rand.Uint32() + 1),
}

var Partners sync.Map

type PartnerTunnel struct {
	sync.RWMutex
	Clients map[uint64]pb.HaClient
}

var Tunnels PartnerTunnel

type NodeInfo struct {
	Id   string
	Addr string
}

const (
	LinkStatus_SHAKEHAND = iota
	LinkStatus_PFAIL
	LinkStatus_FAIL
)

type Link struct {
	Node    NodeInfo
	SendBuf chan pb.Message
	//RecvBuf    chan pb.Message
	MsgContext   context.Context
	Cancel       context.CancelFunc
	Status       uint
	Next         *Link
	Process      MsgProcessor
	IsServerSide bool
}

func (this *Link) SendMsg(msg *pb.Message) {
	if msg == nil {
		return
	}

	this.SendBuf <- *msg
}

type MsgProcessor func(lk *Link, msg *pb.Message) error

func (this *Link) Construct(info NodeInfo, p MsgProcessor) {
	this.Process = p
	this.Node = info
	this.SendBuf = make(chan pb.Message, 10)
	this.MsgContext, this.Cancel = context.WithCancel(context.Background())
}

func (this *Link) RunClientSide(info NodeInfo, client pb.Ha_TwoWayClient) {
	this.IsServerSide = false
	//this.RecvBuf = make(chan pb.Message, 10)

	doRecv := func() error {
		msg, err := client.Recv()
		if err != nil {
			return errors.Wrap(err, "client读操作发生错误")
		}

		//this.RecvBuf <- msg
		if this.Process != nil {
			this.Process(this, msg)
		}

		return nil
	}

	// TODO 启动写协程
	go func(ctx context.Context) {
		var err error
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			default:
				err = doRecv()
				if err != nil {
					break ForLoop
				}
			}
		}

		// TODO 处理link异常
	}(this.MsgContext)

	go func(ctx context.Context) error {
		var err error
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			case msg := <-this.SendBuf:
				err = client.Send(&msg)
				break ForLoop
			}
		}
		// TODO 处理link异常
		return err
	}(this.MsgContext)

}

func (this *Link) RunServerSide(info NodeInfo, ser pb.Ha_TwoWayServer) {
	this.IsServerSide = true

	doRecv := func() error {
		msg, err := ser.Recv()
		if err != nil {
			return errors.Wrap(err, "ser读操作发生错误")
		}

		if this.Process != nil {
			this.Process(this, msg)
		}

		return nil
	}

	// TODO 启动写协程
	go func(ctx context.Context) {
		var err error
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			default:
				err = doRecv()
				if err != nil {
					break ForLoop
				}
			}
		}

		// TODO 处理link异常
	}(this.MsgContext)

	go func(ctx context.Context) error {
		var err error
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			case msg := <-this.SendBuf:
				err = ser.Send(&msg)
				break ForLoop
			}
		}
		// TODO 处理link异常
		return err
	}(this.MsgContext)

}
