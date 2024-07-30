package tundrive

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
)

type TUNPacketMange interface {
	GetPacket() (TUNPacket, error)
	DestructionPacket(TUNPacket) error
}

type TUNDrive struct {
	ctx                context.Context
	lock               sync.Locker
	TunInfo            TUNInfo
	TunFIle            *os.File
	PacketManger       TUNPacketMange
	ReadPacketChannle  chan TUNPacket
	WritePacketChannle chan TUNPacket
	RXPacket           atomic.Uint64
	TXPacket           atomic.Uint64
	RXByte             atomic.Uint64
	TXByte             atomic.Uint64
	TXQueueManage      QueueManage
	DecPakcet          func(TUNPacketMange) TUNPacketMange
	EncPakcet          func(TUNPacketMange) TUNPacketMange
}

func NewDefalutTunDrive(context context.Context) TUNDrive {
	return TUNDrive{
		ctx:                context,
		lock:               &sync.RWMutex{},
		PacketManger:       NewDefalutPacketManger(),
		ReadPacketChannle:  make(chan TUNPacket, 100000),
		WritePacketChannle: make(chan TUNPacket, 100000),
	}
}

func (tun *TUNDrive) NewTun(t TUNInfo) (err error) {
	if err = t.Check(); err != nil {
		return
	}
	if err = t.CreateTunDrive(); err != nil {
		return
	}
	tun.TunInfo = t
	return
}

func (tun *TUNDrive) Start(t TUNInfo) {
	go func(t *TUNDrive) {
		tun.TunFIle = os.NewFile(uintptr(tun.TunInfo.FD), "/dev/net/tun")
		var readPacket TUNPacket

		var err error
		var readNum int
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				readPacket, err = t.PacketManger.GetPacket()
				t.RXPacket.Add(1)
				if err != nil {
					return
				}
				readNum, err = t.TunFIle.Read(readPacket.RawByte)
				t.RXByte.Add(uint64(readNum))
				if err != nil {
					return
				}
				tun.ReadPacketChannle <- readPacket
			}
		}
	}(tun)

	go func(t *TUNDrive) {
		var wirtePack TUNPacket
		var writeNum int
		for {
			select {
			case <-t.ctx.Done():
				return
			case wirtePack = <-t.WritePacketChannle:
				t.TXPacket.Add(1)
				writeNum, _ = t.TunFIle.Write(wirtePack.DecByte)
				t.TXByte.Add(uint64(writeNum))
				t.PacketManger.DestructionPacket(wirtePack)
			}
		}
	}(tun)
}
