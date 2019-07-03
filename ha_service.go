package bulbasaur

import (
	"context"

	"github.com/pkg/errors"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/lj-211/bulbasaur/protocol"
)

type HaServer struct{}

func (s *HaServer) HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (*empty.Empty, error) {
	Logger.Info("SER-收到心跳包: ", req.Id)
	return &empty.Empty{}, nil
}

func (s *HaServer) GetNodeList(ctx context.Context, req *pb.GetNodeListReq) (*pb.GetNodeListRes, error) {
	return nil, errors.New("not implemented")
}

func (s *HaServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	// 1. 检查节点是否存在
	if v, ok := bulbasaur.Partners[req.Id]; ok {
		v.LastActive = time.Now()
		common.Log.Warningf("parterner-%+v已注册，重复收到注册消息", req.Id)
	} else {
		bulbasaur.Partners[req.Id] = &bulbasaur.NodeInfo{
			NodeId:     req.Id,
			Addr:       req.Addr,
			LastActive: time.Now(),
			//Role: 1,
		}
		common.Log.Infof("新的parterner-%+v加入", req.Id)
	}
	return nil, errors.New("not implemented")
}
