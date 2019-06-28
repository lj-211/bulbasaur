package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lj-211/grpcwrapper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/lj-211/bulbasaur"
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
			res, cerr := haClient.HeartBeat(context.Background(), &pb.HeartBeatReq{
				Id: 2,
			})
			if cerr != nil {
				fmt.Println("心跳发生错误: ", cerr.Error())
			}
			time.Sleep(time.Second * 3)
		}
	}()

	return nil
}

const serverAddr string = "0.0.0.0:8000"
const selfAddr string = "0.0.0.0:8000"

func main() {
	fmt.Println("start running...")
	var exitErr error
	for ok := true; ok; ok = false {
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
