package bulbasaur

import (
	"math/rand"
	//"sync"
	"time"

	pb "github.com/lj-211/bulbasaur/protocol"
)

type NodeInfo struct {
	NodeId     uint64
	Addr       string // 暂时没有服务发现，使用ip:port
	LastActive time.Time
	Role       pb.Role
}

type NodeMap map[uint64]*NodeInfo

var Info *NodeInfo = &NodeInfo{
	NodeId: uint64(rand.Uint32() + 1),
}

var Partners NodeMap = make(NodeMap)
