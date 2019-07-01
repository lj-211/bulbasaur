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
	example "github.com/lj-211/bulbasaur/example"
	pb "github.com/lj-211/bulbasaur/protocol"
)

func _registProto(e *grpcwrapper.Engine) {
	pb.RegisterHaServer(e.GetRawServer(), new(bulbasaur.HaServer))
}

func _registMiddleware(e *grpcwrapper.Engine) {
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

	go func() {
		for {
			_, cerr := haClient.HeartBeat(context.Background(), &pb.HeartBeatReq{
				Id: 2,
			})
			if cerr != nil {
				if ec, ok := ecode.Cause(cerr).(ecode.Codes); ok {
					fmt.Println("CLI-心跳返回标准错误码: ", ec.Code(), ec.Message())
				} else {
					fmt.Println("CLI-心跳发生错误: ", cerr.Error())
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()

	return nil
}

var serverAddr string
var selfAddr string

func paramParse() error {
	flag.StringVar(&serverAddr, "remote", "", "远端地址")
	flag.StringVar(&selfAddr, "local", "", "本地监听地址")

	flag.Parse()

	if selfAddr == "" {
		return errors.New("本地监听地址未设置")
	}

	if serverAddr == "" {
		return errors.New("伙伴地址至少需要设定一个")
	}

	return nil
}

func main() {
	fmt.Println("start running...")
	var exitErr error
	for ok := true; ok; ok = false {
		if err := paramParse(); err != nil {
			exitErr = errors.Wrap(err, "启动参数设定非法")
			break
		}
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
