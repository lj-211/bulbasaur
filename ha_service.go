package bulbasaur

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	//"github.com/lj-211/common/ecode"
	//"github.com/lj-211/grpcwrapper"
	"github.com/pkg/errors"
	//"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	//empty "github.com/golang/protobuf/ptypes/empty"
	common "github.com/lj-211/bulbasaur/example"
	pb "github.com/lj-211/bulbasaur/protocol"
)

type HaServer struct{}

func (this *HaServer) TwoWay(stream pb.Ha_TwoWayServer) error {
	ctx := stream.Context()
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return nil
	}

	common.Log.Infof("peer发送消息: %+v", *pr)

	lk := &Link{}
	info := NodeInfo{
		// 握手阶段id先留空
		Addr: pr.Addr.String(),
	}
	lk.Construct(info, processMsg)
	AddLink(&MySelf, lk)

	lk.RunServerSide(stream)

	return nil
}

func PingProc(ctx context.Context, lk *Link, msg *pb.Message) error {
	if lk == nil {
		return errors.New("连接为空")
	}
	if msg == nil {
		return errors.New("消息为空")
	}

	ping := &pb.Ping{}
	if err := proto.Unmarshal(msg.Data, ping); err != nil {
		return errors.Wrapf(err, "解码消息出错")
	}

	common.Log.Infof("伙伴发来心跳消息")

	if lk.Status == LinkStatus_SHAKEHAND {
		lk.Status = LinkStatus_ACTIVE
	}
	lk.LastActive = time.Now()
	lk.Node.Id = ping.Id

	common.Log.Infof("%s响应ping,完成握手", ping.Id)

	pmsg := &pb.Pong{
		Id: MySelf.Id,
	}
	sendMsg(lk, MTypePong, pmsg)

	return nil
}

func PongProc(ctx context.Context, lk *Link, msg *pb.Message) error {
	if lk == nil {
		return errors.New("连接为空")
	}
	if msg == nil {
		return errors.New("消息为空")
	}

	pong := &pb.Pong{}
	if err := proto.Unmarshal(msg.Data, pong); err != nil {
		return errors.Wrapf(err, "解码消息出错")
	}

	if lk.Status == LinkStatus_SHAKEHAND {
		common.Log.Infof("设置link状态为存活")
		lk.Status = LinkStatus_ACTIVE
	}
	lk.LastActive = time.Now()
	lk.Node.Id = pong.Id

	common.Log.Infof("%s响应pong,完成握手", pong.Id)

	return nil
}

func init() {
	registerMsgProcess(MTypePing, PingProc)
	registerMsgProcess(MTypePong, PongProc)
}
