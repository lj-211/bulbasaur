package bulbasaur

import (
	"context"
	"fmt"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/lj-211/bulbasaur/protocol"
)

type HaServer struct{}

func (s *HaServer) HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (*empty.Empty, error) {
	fmt.Println("收到心跳包: ", req.Id)
	return &empty.Empty{}, nil
}
