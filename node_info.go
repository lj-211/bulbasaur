package bulbasaur

import (
	"math/rand"
	"sync"
	"time"

	pb "github.com/lj-211/bulbasaur/protocol"
)

type PartnerInfo struct {
	PartnerId  uint64
	Addr       string // 暂时没有服务发现，使用ip:port
	LastActive time.Time
	Role       pb.Role
}

var Info *PartnerInfo = &PartnerInfo{
	PartnerId: uint64(rand.Uint32() + 1),
}

var Partners sync.Map

type PartnerTunnel struct {
	sync.RWMutex
	Clients map[uint64]pb.HaClient
}

var Tunnels PartnerTunnel
