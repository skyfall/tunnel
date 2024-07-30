package tundrive

import "sync/atomic"

type RateInfo struct {
	Packet atomic.Uint64
	Byte   atomic.Uint64
}

type QueueManage interface {
	SetTunPakcet(Packet TUNPacket) (err error)
}
