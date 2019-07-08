package bulbasaur

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	pb "github.com/lj-211/bulbasaur/protocol"
)

const (
	MtypeHeartBeat = iota
)

type MsgProcess func(context.Context, uint64, []byte, pb.Ha_TwoWayServer) error

func processMsg(ctx context.Context, id uint64, msg *pb.Message, tw pb.Ha_TwoWayServer) error {
	if id == 0 {
		return errors.New("没有合法的伙伴id")
	}

	if p, ok := MsgMap[msg.Mtype]; ok {
		return p(ctx, id, msg.Data, tw)
	} else {
		return errors.Errorf("位置消息%d类型", msg.Mtype)
	}
	return nil
}

var MsgMap map[uint32]MsgProcess = make(map[uint32]MsgProcess)

func registerMsgProcess(mtype uint32, p MsgProcess) error {
	if _, ok := MsgMap[mtype]; ok {
		return errors.Errorf("消息%d处理已存在", mtype)
	}
	MsgMap[mtype] = p

	return nil
}

func sendMsg(tw pb.Ha_TwoWayServer, mtype int, msg proto.Message) {
	buf, _ := proto.Marshal(msg)
	tw.Send(&pb.Message{
		Mtype: uint32(mtype),
		Data:  buf,
	})
}
