package bulbasaur

import (
	"context"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lj-211/common/ecode"
	"github.com/lj-211/grpcwrapper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	empty "github.com/golang/protobuf/ptypes/empty"
	common "github.com/lj-211/bulbasaur/example"
	pb "github.com/lj-211/bulbasaur/protocol"
)

type HaServer struct{}

func (s *HaServer) HeartBeat(ctx context.Context, req *pb.HeartBeatReq) (*empty.Empty, error) {
	common.Log.Infof("收到partner-%v心跳", req.Id)
	if _, ok := Partners.Load(req.Id); ok {
		common.Log.Infof("新的partner-%v", req.Id)
		Partners.Store(req.Id, &PartnerInfo{
			PartnerId:  req.Id,
			Addr:       req.Addr,
			LastActive: time.Now(),
		})
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

		clientOpt := grpcwrapper.DefaultClient()
		client, derr := clientOpt.DialContext(context.Background(), req.Addr, grpc.WithBlock())
		if derr != nil {
			return nil, ecode.Errorf(ecode.Code(pb.ErrCode_ConnectFail), "尝试连接失败，请重新注册")
		}
		haClient := pb.NewHaClient(client)

		Tunnels.Lock()
		defer Tunnels.Unlock()
		Tunnels.Clients[req.Id] = haClient
		common.Log.Info("已和%+v连接成功")
	}
	return &pb.RegisterRes{}, nil
}

func (this *HaServer) TwoWay(stream pb.Ha_TwoWayServer) error {
	ctx := stream.Context()
	pr, ok := peer.FromContext(ctx)
	if ok {
		common.Log.Infof("peer发送消息: %+v", *pr)
	}

	var pid uint64
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			common.Log.Info("伙伴连接断开")
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "处理双向消息出错")
		}

		if msg.Mtype == MtypeHeartBeat {
			hb := &pb.HeartBeatReq{}
			if err := proto.Unmarshal(msg.Data, hb); err != nil {
				common.Log.Errorf("解码消息出错: %s", err.Error())
				continue
			}
			pid = hb.Id
		}

		if perr := processMsg(context.Background(), pid, msg, stream); perr != nil {
			common.Log.Errorf("处理消息时发生错误: %s", perr.Error())
		}
	}
}

func HeartBeat(ctx context.Context, pid uint64, data []byte, tw pb.Ha_TwoWayServer) error {
	hb := &pb.HeartBeatReq{}
	if err := proto.Unmarshal(data, hb); err != nil {
		return errors.Wrapf(err, "解码消息出错")
	}

	common.Log.Infof("伙伴%d发来心跳消息", hb.Id)
	sendMsg(tw, MtypeHeartBeat, hb)

	return nil
}

func init() {
	registerMsgProcess(MtypeHeartBeat, HeartBeat)
}
