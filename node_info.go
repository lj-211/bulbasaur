package bulbasaur

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
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
		ni.LinkTail.Pre = ni.LinkHead
		common.Log.Info("init link list")
		return
	}

	old := ni.LinkTail
	lk.Pre = old
	lk.Next = nil
	old.Next = lk
	ni.LinkTail = lk
	common.Log.Info("add link list")
	node := ni.LinkHead
	common.Log.Info(node)
	for node != nil {
		common.Log.Infof("节点: %s", node.Node.Id)
		node = node.Next
	}
	common.Log.Info("******************")
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
	away         chan bool
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
	this.away = make(chan bool)

	go func() {
		<-this.away
		common.Log.Infof("%s丢失", this.Node.Addr)
		this.Cancel()
		this.Status = LinkStatus_FAIL

		// clear channel
		for range this.away {
		}
	}()
}

func (this *Link) Goaway() {
	this.away <- true
	common.Log.Info("goaway")
}

func (this *Link) RunClientSide(client pb.Ha_TwoWayClient) {
	this.IsServerSide = false

	rctx, _ := context.WithCancel(this.MsgContext)
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
					this.Goaway()
					return
				}
				if terr != nil {
					common.Log.Infof("-client recv err %s", terr.Error())
					this.Goaway()
					return
				}
				if msg != nil && this.Process != nil {
					common.Log.Infof("client recv msg %+v", msg)
					this.Process(context.TODO(), this, msg)
				}
				time.Sleep(time.Second * 2)
			}
		}
	}(rctx)

	sctx, _ := context.WithCancel(this.MsgContext)
	go func(ctx context.Context) {
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			case msg := <-this.SendBuf:
				err := client.Send(&msg)
				if err != nil {
					common.Log.Infof("client side error %s", err.Error())
					this.Goaway()
				}
			}
		}
	}(sctx)

}

func (this *Link) RunServerSide(ser pb.Ha_TwoWayServer) {
	this.IsServerSide = true

	// send
	sctx, _ := context.WithCancel(this.MsgContext)
	go func(ctx context.Context) {
	ForLoop:
		for {
			select {
			case <-ctx.Done():
				break ForLoop
			case msg := <-this.SendBuf:
				common.Log.Info("服务器发送消息")
				err := ser.Send(&msg)
				if err != nil {
					this.Goaway()
				}
				break ForLoop
			}
		}
	}(sctx)

ForLoop:
	for {
		select {
		case <-this.MsgContext.Done():
			break ForLoop
		default:
			msg, terr := ser.Recv()
			if terr == io.EOF {
				common.Log.Error("ser收取消息错误: ", terr.Error())
				this.Goaway()
				break ForLoop
			}
			if terr != nil {
				break ForLoop
			}
			if msg != nil && this.Process != nil {
				common.Log.Infof("server recv msg %+v", msg)
				this.Process(context.TODO(), this, msg)
			}
			time.Sleep(time.Second * 2)
		}
	}

	// TODO 处理link异常
	this.Goaway()
}

func NodeOk(id string) {
	LinkLock.RLock()
	LinkLock.RUnlock()

	node := MySelf.LinkHead
	for node != nil {
		if node.Node.Id == id {
			bmsg := &pb.Ping{
				Id: MySelf.Id,
			}
			buf, _ := proto.Marshal(bmsg)
			node.SendMsg(&pb.Message{
				Mtype: MTypePing,
				Data:  buf,
			})
			break
		}

		node = node.Next
	}
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

const MaxLostTime time.Duration = time.Second * 5

func NodeCron() {
	LinkLock.Lock()
	LinkLock.Unlock()

	node := MySelf.LinkHead
	for node != nil {
		next := node.Next
		for ok := true; ok; ok = false {
			if node.Status == LinkStatus_PFAIL {
				now := time.Now()
				if now.Sub(node.LastActive) > 2*MaxLostTime {
					node.Status = LinkStatus_FAIL
				}
				break
			}
			// 1. 删除已丢失状态的link
			if node.Status == LinkStatus_FAIL {
				if node == MySelf.LinkHead {
					MySelf.LinkHead = next
				}
				if next == nil {
					node.Pre.Next = nil
					MySelf.LinkTail = node.Pre
				} else {
					node.Pre.Next = node.Next
					node.Next.Pre = node.Pre
				}

				common.Log.Infof("已丢失节点%s被剔除连接列表", node.Node.Id)
				break
			}

			// 2. 检查客户端link active时间，标志状态为PFAIL，并且发送IsOK消息
			if node.Status == LinkStatus_ACTIVE {
				now := time.Now()
				diff := now.Sub(node.LastActive)
				if diff > MaxLostTime && node.IsServerSide {
					node.Status = LinkStatus_PFAIL
					// 发送IsOk查询
					bmsg := &pb.IsOk{
						Id: MySelf.Id,
					}
					buf, _ := proto.Marshal(bmsg)
					node.SendMsg(&pb.Message{
						Mtype: MTypeIsOk,
						Data:  buf,
					})
				}

				if diff > MaxLostTime/3 && !node.IsServerSide {
					bmsg := &pb.Ping{
						Id: MySelf.Id,
					}
					buf, _ := proto.Marshal(bmsg)
					node.SendMsg(&pb.Message{
						Mtype: MTypePing,
						Data:  buf,
					})
				}
			}
		}

		node = next
	}
}
