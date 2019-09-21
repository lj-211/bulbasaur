package bulbasaur

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	pb "github.com/lj-211/bulbasaur/protocol"
)

const (
	MTypeHeartBeat = iota
	MTypePing
	MTypePong
	MTypeIsOk
)

type MsgProcessor func(context.Context, *Link, *pb.Message) error

func processMsg(ctx context.Context, lk *Link, msg *pb.Message) error {
	if lk == nil {
		return nil
	}

	if p, ok := MsgMap[msg.Mtype]; ok {
		return p(ctx, lk, msg)
	} else {
		return errors.Errorf("未知消息%d类型", msg.Mtype)
	}

	return nil
}

var MsgMap map[uint32]MsgProcessor = make(map[uint32]MsgProcessor)

func registerMsgProcess(mtype uint32, p MsgProcessor) error {
	if _, ok := MsgMap[mtype]; ok {
		return errors.Errorf("消息%d处理已存在", mtype)
	}
	MsgMap[mtype] = p

	return nil
}

func sendMsg(lk *Link, mtype int, msg proto.Message) {
	if lk == nil {
		return
	}
	buf, _ := proto.Marshal(msg)
	lk.SendMsg(&pb.Message{
		Mtype: uint32(mtype),
		Data:  buf,
	})
}
