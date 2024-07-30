package tundrive

import (
	"context"
	"sync/atomic"
)

const IPPROTO_DEFALUT = 0x0

type QueueProtocolManage struct {
	ctx   context.Context
	Queue map[int]QueueProtocol
}

type QueueProtocol struct {
	Rate             RateInfo
	ctx              context.Context
	Packet           chan TUNPacket
	PacketEncDncFunc func(TUNPacket)
	SendPacketFunc   func(TUNPacket)
}

func (protocol QueueProtocol) Start() {
	go func(q QueueProtocol) {
		var packet TUNPacket
		for {
			select {
			case <-q.ctx.Done():
				return
			case packet = <-q.Packet:
				q.PacketEncDncFunc(packet)
				if q.SendPacketFunc != nil {
					q.SendPacketFunc(packet)
				}
			}
		}
	}(protocol)
}

func NewQueueProtocol(ctx context.Context, packetEncDncFunc func(TUNPacket), SendPacketFunc func(TUNPacket)) (q QueueProtocolManage) {
	IPPROTO_TCP := 0x6
	IPPROTO_UDP := 0x11
	IPPROTO_ICMP := 0x1
	childCtx, _ := context.WithCancel(ctx)
	q = QueueProtocolManage{
		ctx:   ctx,
		Queue: make(map[int]QueueProtocol, 5),
		// QueueStatus: make(map[int]QueueStatus, 5),
	}

	q.Queue[IPPROTO_DEFALUT] = QueueProtocol{
		ctx:              childCtx,
		Packet:           make(chan TUNPacket, 10000),
		PacketEncDncFunc: packetEncDncFunc,
		SendPacketFunc:   SendPacketFunc,
		Rate:             RateInfo{Packet: atomic.Uint64{}, Byte: atomic.Uint64{}},
	}
	q.Queue[IPPROTO_DEFALUT].Start()

	q.Queue[IPPROTO_TCP] = QueueProtocol{
		ctx:              childCtx,
		Packet:           make(chan TUNPacket, 10000),
		PacketEncDncFunc: packetEncDncFunc,
		SendPacketFunc:   SendPacketFunc,
		Rate:             RateInfo{Packet: atomic.Uint64{}, Byte: atomic.Uint64{}},
	}
	q.Queue[IPPROTO_TCP].Start()

	q.Queue[IPPROTO_UDP] = QueueProtocol{
		ctx:              childCtx,
		Packet:           make(chan TUNPacket, 10000),
		PacketEncDncFunc: packetEncDncFunc,
		SendPacketFunc:   SendPacketFunc,
		Rate:             RateInfo{Packet: atomic.Uint64{}, Byte: atomic.Uint64{}},
	}
	q.Queue[IPPROTO_UDP].Start()

	q.Queue[IPPROTO_ICMP] = QueueProtocol{
		ctx:              childCtx,
		Packet:           make(chan TUNPacket, 10000),
		PacketEncDncFunc: packetEncDncFunc,
		SendPacketFunc:   SendPacketFunc,
		Rate:             RateInfo{Packet: atomic.Uint64{}, Byte: atomic.Uint64{}},
	}
	q.Queue[IPPROTO_ICMP].Start()

	return
}

func (q QueueProtocolManage) SetTunPakcet(Packet TUNPacket) (err error) {
	var queue QueueProtocol
	var ok bool
	queue, ok = q.Queue[Packet.Protocol]
	if ok {
		queue = q.Queue[IPPROTO_DEFALUT]
	}
	queue.Packet <- Packet
	queue.Rate.Packet.Add(1)
	queue.Rate.Byte.Add(uint64(len(Packet.RawByte)))
	return
}
