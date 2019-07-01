package bulbasaur

import (
	"time"
)

type NodeInfo struct {
	NodeId     uint32
	Addr       string // 暂时没有服务发现，使用ip:port
	LastActive time.Time
	Role       int32
}
