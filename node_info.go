package bulbasaur

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/lj-211/grpcwrapper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	common "github.com/lj-211/bulbasaur/example"
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
	Id       string
	Addr     string
	LinkHead *Link
	LinkTail *Link
}

var MySelf NodeInfo
var LinkLock sync.RWMutex

func PrintLinkList() {
	LinkLock.RLock()
	defer LinkLock.RUnlock()

	common.Log.Info("节点列表: ")
	common.Log.Info("====================")
	node := MySelf.LinkHead
	for node != nil {
		common.Log.Infof("节点: %s 地址: %s, 状态: %s", node.Node.Id,
			node.Node.Addr, getLinkStatusStr(node.Status))
		node = node.Next
	}
	common.Log.Info("====================")
}

func AddLink(ni *NodeInfo, lk *Link) {
	if lk == nil {
		return
	}

	LinkLock.Lock()
	defer LinkLock.Unlock()

	if ni.LinkHead == nil {
		ni.LinkHead = lk
		ni.LinkTail = lk
		return
	}

	lk.Pre = ni.LinkTail
	lk.Next = nil
	ni.LinkTail = lk
}

const (
	LinkStatus_UNKNOWN = iota
	LinkStatus_SHAKEHAND
	LinkStatus_PFAIL
	LinkStatus_FAIL
	LinkStatus_ACTIVE
)

func getLinkStatusStr(status uint) string {
	s_str := ""
	switch status {
	case LinkStatus_UNKNOWN:
		s_str = "未知"
	case LinkStatus_SHAKEHAND:
		s_str = "握手"
	case LinkStatus_PFAIL:
		s_str = "丢失待确认"
	case LinkStatus_FAIL:
		s_str = "丢失"
	case LinkStatus_ACTIVE:
		s_str = "存活"
	default:
		s_str = "非法状态"
	}

	return s_str
}

type Link struct {
	Node    NodeInfo
	SendBuf chan pb.Message
	//RecvBuf    chan pb.Message
	MsgContext   context.Context
	Cancel       context.CancelFunc
	Status       uint
	Pre          *Link
	Next         *Link
	Process      MsgProcessor
	IsServerSide bool
	LastActive   time.Time
}

func (this *Link) SendMsg(msg *pb.Message) {
	if msg == nil {
		return
	}

	this.SendBuf <- *msg
}

func (this *Link) Construct(info NodeInfo, p MsgProcessor) {
	this.Process = p
	this.Node = info
	this.Status = LinkStatus_SHAKEHAND
	this.SendBuf = make(chan pb.Message, 10)
	this.MsgContext, this.Cancel = context.WithCancel(context.Background())
	this.LastActive = time.Now()
}

func (this *Link) RunClientSide(client pb.Ha_TwoWayClient) {
	this.IsServerSide = false

	go func(ctx context.Context) {
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			default:
				msg, terr := client.Recv()
				if terr == io.EOF {
					common.Log.Error("收取消息错误: ", terr.Error())
					break ForLoop
				}
				if terr != nil {
					common.Log.Infof("client recv err %s", terr.Error())
					break
				}
				if msg != nil && this.Process != nil {
					common.Log.Infof("client recv msg %+v", msg)
					this.Process(context.TODO(), this, msg)
				}
				time.Sleep(time.Second * 2)
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

func (this *Link) RunServerSide(ser pb.Ha_TwoWayServer) {
	this.IsServerSide = true

	// send
	go func(ctx context.Context) error {
		var err error
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			case msg := <-this.SendBuf:
				common.Log.Info("服务器发送消息")
				err = ser.Send(&msg)
				break ForLoop
			}
		}
		// TODO 处理link异常
		return err
	}(this.MsgContext)

ForLoop:
	for {
		select {
		case <-this.MsgContext.Done():
			break ForLoop
		default:
			msg, terr := ser.Recv()
			if terr == io.EOF {
				common.Log.Error("ser收取消息错误: ", terr.Error())
				break ForLoop
			}
			if terr != nil {
				break
			}
			if msg != nil && this.Process != nil {
				common.Log.Infof("server recv msg %+v", msg)
				this.Process(context.TODO(), this, msg)
			}
			time.Sleep(time.Second * 2)
		}
	}

	// TODO 处理link异常
}

func ConnectNode(id int, addr string) (*Link, error) {
	clientOpt := grpcwrapper.DefaultClient()
	client, derr := clientOpt.DialContext(context.Background(), addr, grpc.WithBlock())
	if derr != nil {
		return nil, errors.Wrap(derr, "连接服务器失败")
	}
	haClient := pb.NewHaClient(client)

	tw, twErr := haClient.TwoWay(context.Background())
	if twErr != nil {
		return nil, errors.Wrap(twErr, "连接节点失败")
	}

	lk := &Link{}
	info := NodeInfo{
		// 握手阶段id先留空
		Addr: addr,
	}
	lk.Construct(info, processMsg)
	AddLink(&MySelf, lk)

	go lk.RunClientSide(tw)

	return lk, nil
}

func NodeCron() {
	LinkLock.Lock()
	LinkLock.Unlock()

	// 1. 检查状态
	node := MySelf.LinkHead
	for node != nil {
		node = node.Next
	}
}
