package bulbasaur

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/lj-211/bulbasaur/protocol"
)

const (
	MtypeHeartBeat = iota
)

type MsgProcess func(context.Context, uint64, []byte, grpc.ServerStream) error

func processMsg(ctx context.Context, id uint64, msg *pb.Message, stream grpc.ServerStream) error {
	if id == 0 {
		return errors.New("没有合法的伙伴id")
	}

	if p, ok := MsgMap[msg.Mtype]; ok {
		return p(ctx, id, msg.Data, stream)
	} else {
		return errors.Errorf("位置消息%d类型", msg.Mtype)
	}
	return nil
}

var MsgMap map[uint32]MsgProcess

func registerMsgProcess(mtype uint32, p MsgProcess) error {
	if _, ok := MsgMap[mtype]; ok {
		return errors.Errorf("消息%d处理已存在", mtype)
	}
	MsgMap[mtype] = p

	return nil
}
