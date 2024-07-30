package tundrive

import (
	"context"
	"encoding/binary"
	"sync/atomic"
)

type QueueHashManage struct {
	ctx   context.Context
	num   int
	Queue []QueueHash
}

type QueueHash struct {
	ctx              context.Context
	Rate             RateInfo
	Packet           chan TUNPacket
	PacketEncDncFunc func(TUNPacket)
	SendPacket       func(TUNPacket)
}

func (hash QueueHash) Start() {
	go func(q QueueHash) {
		var packet TUNPacket
		for {
			select {
			case <-q.ctx.Done():
				return
			case packet = <-q.Packet:
				q.PacketEncDncFunc(packet)
				if q.SendPacket != nil {
					q.SendPacket(packet)
				}

			}
		}
	}(hash)
}

func NewQueueHash(ctx context.Context, num int, packetEncDncFunc func(TUNPacket), SendPacketFunc func(TUNPacket)) (q QueueHashManage) {
	q = QueueHashManage{}
	q.ctx = ctx
	q.Queue = make([]QueueHash, num)
	q.num = num
	childCtx, _ := context.WithCancel(ctx)
	for index := 0; index < q.num; index++ {
		q.Queue[index] = QueueHash{
			ctx:              childCtx,
			Packet:           make(chan TUNPacket, 10000),
			PacketEncDncFunc: packetEncDncFunc,
			SendPacket:       SendPacketFunc,
			Rate:             RateInfo{Packet: atomic.Uint64{}, Byte: atomic.Uint64{}},
		}
		q.Queue[index].Start()
	}
	return
}

func (q QueueHashManage) SetTunPakcet(Packet TUNPacket) (err error) {
	var value uint64
	value = 0
	value += uint64(Packet.DesPort)
	value += uint64(Packet.Protocol)
	value += binary.LittleEndian.Uint64(Packet.DesIP)
	index := value % uint64(q.num)

	q.Queue[index].Rate.Packet.Add(1)
	q.Queue[index].Rate.Byte.Add(uint64(len(Packet.RawByte)))
	q.Queue[index].Packet <- Packet
	return
}
