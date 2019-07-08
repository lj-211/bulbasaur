package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lj-211/common/ecode"
	"github.com/lj-211/grpcwrapper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/lj-211/bulbasaur"
	common "github.com/lj-211/bulbasaur/example"
	pb "github.com/lj-211/bulbasaur/protocol"
)

func _registProto(e *grpcwrapper.Engine) {
	pb.RegisterHaServer(e.GetRawServer(), new(bulbasaur.HaServer))
}

func _registMiddleware(e *grpcwrapper.Engine) {
}

func initSelfPartnerInfo() {
	bulbasaur.Info.Addr = selfAddr
	if partnerId != 0 {
		bulbasaur.Info.PartnerId = uint64(partnerId)
	}

	bulbasaur.Tunnels.Clients = make(map[uint64]pb.HaClient)
}

func startServer(addr string) error {
	engine := grpcwrapper.Default()
	_registMiddleware(engine)
	if err := engine.InitServer(); err != nil {
		return errors.Wrap(err, "启动服务失败")
	}
	_registProto(engine)
	go engine.Run(addr)
	return nil
}

func connectServer(addr string) error {
	clientOpt := grpcwrapper.DefaultClient()
	client, derr := clientOpt.DialContext(context.Background(), addr, grpc.WithBlock())
	if derr != nil {
		return errors.Wrap(derr, "连接服务器失败")
	}
	haClient := pb.NewHaClient(client)

	tw, twErr := haClient.TwoWay(context.Background())
	if twErr != nil {
		common.Log.Errorf("创建双向客户端失败: %s", twErr.Error())
	} else {
		go func() {
			for {
				time.Sleep(time.Second * 2)
				msg, merr := tw.Recv()
				if merr == io.EOF {
					common.Log.Infof("伙伴%d断开连接", 123456)
					return
				}
				if merr != nil {
					common.Log.Infof("读数据发生错误 %s", merr.Error())
					break
				} else {
					common.Log.Info("收到消息 %s", strconv.FormatUint(uint64(msg.Mtype), 10))
				}

				bmsg := &pb.HeartBeatReq{
					Id: bulbasaur.Info.PartnerId,
				}
				buf, _ := proto.Marshal(bmsg)
				tw.Send(&pb.Message{
					Mtype: bulbasaur.MtypeHeartBeat,
					Data:  buf,
				})
				common.Log.Info("发送心跳消息")

				common.Log.Info("two way client tick")
			}
		}()
	}

	for false {
		_, err := haClient.Register(context.Background(), &pb.RegisterReq{
			Id:   uint64(bulbasaur.Info.PartnerId),
			Addr: bulbasaur.Info.Addr,
		})
		if err != nil {
			ec, _ := ecode.Cause(err).(ecode.Codes)
			common.Log.Info("连接远端发生错误: ", ec.Message())
			time.Sleep(time.Second * 2)
			continue
		}

		time.Sleep(time.Second)

		break
	}

	common.Log.Infof("开始和%s心跳", addr)

	go func() {
		for {
			_, cerr := haClient.HeartBeat(context.Background(), &pb.HeartBeatReq{
				Id: bulbasaur.Info.PartnerId,
			})
			if cerr != nil {
				if ec, ok := ecode.Cause(cerr).(ecode.Codes); ok {
					common.Log.Info("CLI-心跳返回标准错误码: ", ec.Code(), ec.Message())
				} else {
					common.Log.Info("CLI-心跳发生错误: ", cerr.Error())
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()

	return nil
}

var serverAddr string
var selfAddr string
var partnerId uint

func paramParse() error {
	flag.StringVar(&serverAddr, "remote", "", "远端地址")
	flag.StringVar(&selfAddr, "local", "", "本地监听地址")
	flag.UintVar(&partnerId, "id", 0, "节点id(大于0)")

	flag.Parse()

	if selfAddr == "" {
		return errors.New("本地监听地址未设置")
	}

	if partnerId == 0 {
		common.Log.Infof("没有设置节点id,将使用系统随机id %d", bulbasaur.Info.PartnerId)
	}

	return nil
}

func main() {
	fmt.Println("start running...")
	var exitErr error
	for ok := true; ok; ok = false {
		if err := common.InitGoLoggingStdout("base"); err != nil {
			exitErr = errors.Wrap(err, "初始化stdout日志失败")
			break
		}
		bulbasaur.Logger.Log = common.Log

		if err := paramParse(); err != nil {
			exitErr = errors.Wrap(err, "启动参数设定非法")
			break
		}
		common.Log.Info("日志初始化成功")

		if err := startServer(selfAddr); err != nil {
			exitErr = errors.Wrap(err, "启动服务失败")
			break
		}

		initSelfPartnerInfo()

		if serverAddr != "" {
			if err := connectServer(serverAddr); err != nil {
				exitErr = errors.Wrap(err, "连接服务失败")
				break
			}
			bulbasaur.Info.Role = pb.Role_Follower
		} else {
			bulbasaur.Info.Role = pb.Role_Leader
			common.Log.Info("作为主节点运行")
		}

		for {
			time.Sleep(time.Second * 5)
			switch bulbasaur.Info.Role {
			case pb.Role_Leader:
				common.Log.Info("我是Leader")
			case pb.Role_Follower:
				common.Log.Info("我是Follower")
			case pb.Role_Candidate:
				common.Log.Info("我是Candidate")
			default:
				common.Log.Error("角色异常")
			}

			if bulbasaur.Info.Role == pb.Role_Leader {
				common.Log.Info("partners:  start")
				bulbasaur.Partners.Range(func(k, v interface{}) bool {
					pinfo := v.(*bulbasaur.PartnerInfo)
					common.Log.Infof("partners: %v %+v", k.(uint64), *pinfo)
					return true
				})
				common.Log.Info("partners: end")
			}
		}
	}

	fmt.Println("服务退出: ", exitErr.Error())
}
