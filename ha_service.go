package bulbasaur

import (
	"context"
	"time"

	"github.com/lj-211/common/ecode"
	"github.com/pkg/errors"

	empty "github.com/golang/protobuf/ptypes/empty"
	common "github.com/lj-211/bulbasaur/example"
	pb "github.com/lj-211/bulbasaur/protocol"
)

type HaServer struct{}

func (s *HaServer) HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (*empty.Empty, error) {
	if _, ok := Partners.Load(req.Id); ok {
		Partners.Store(req.Id, &PartnerInfo{
			PartnerId:  req.Id,
			Addr:       req.Addr,
			LastActive: time.Now(),
		})
		common.Log.Info("收到partner-%v心跳", req.Id)
	} else {
		common.Log.Warningf("未注册的partner-%v的心跳被忽略", req.Id)
	}
	return &empty.Empty{}, nil
}

func (s *HaServer) GetNodeList(ctx context.Context, req *pb.GetNodeListReq) (*pb.GetNodeListRes, error) {
	return nil, errors.New("not implemented")
}

func (s *HaServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	common.Log.Infof("收到注册消息: %+v", req.Id)
	// 1. 检查节点是否存在
	if _, ok := Partners.Load(req.Id); ok {
		common.Log.Warningf("parterner-%+v已注册，重复收到注册消息", req.Id)
		return nil, ecode.Errorf(ecode.Code(pb.ErrCode_PartnerHasRegistered),
			"parterner-%+v已注册", req.Id)
	} else {
		Partners.Store(req.Id, &PartnerInfo{
			PartnerId:  req.Id,
			Addr:       req.Addr,
			LastActive: time.Now(),
			//Role: 1,
		})
		common.Log.Infof("新的parterner-%+v加入", req.Id)
	}
	return &pb.RegisterRes{}, nil
}
