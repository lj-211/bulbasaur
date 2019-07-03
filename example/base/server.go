package main

import (
	"context"
	"flag"
	"fmt"
	"time"

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

func initSelfNodeInfo() {
	bulbasaur.Info.Addr = selfAddr
	bulbasaur.Info.Role = pb.Role_Leader
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

	for {
		_, err := haClient.Register(context.Background(), &pb.RegisterReq{
			Id:   uint64(bulbasaur.Info.NodeId),
			Addr: bulbasaur.Info.Addr,
		})
		if err != nil {
			common.Log.Info("连接远端成功")
		}

		time.Sleep(time.Second)

		break
	}

	common.Log.Infof("开始和%s心跳", addr)
	go func() {
		for {
			_, cerr := haClient.HeartBeat(context.Background(), &pb.HeartBeatReq{
				Id: 2,
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
var nodeId uint

func paramParse() error {
	flag.StringVar(&serverAddr, "remote", "", "远端地址")
	flag.StringVar(&selfAddr, "local", "", "本地监听地址")
	flag.UintVar(&nodeId, "id", 0, "节点id(大于0)")

	flag.Parse()

	if selfAddr == "" {
		return errors.New("本地监听地址未设置")
	}

	if serverAddr == "" {
		return errors.New("伙伴地址至少需要设定一个")
	}

	if nodeId == 0 {
		common.Log.Infof("没有设置节点id,将使用系统随机id %d", bulbasaur.Info.NodeId)
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

		if err := connectServer(serverAddr); err != nil {
			exitErr = errors.Wrap(err, "连接服务失败")
			break
		}

		for {
			time.Sleep(time.Second * 1)
		}
	}

	fmt.Println("服务退出: ", exitErr.Error())
}
