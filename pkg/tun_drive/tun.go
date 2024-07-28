package tundrive

import (
	"context"
	"sync"
)

type TUNReadPacketManage interface {
	GetPacket() ([]byte, error)
	DestructionPacket([]byte) error
}

type TUNWritePacketManage interface {
	GetPacket() ([]byte, error)
	DestructionPacket([]byte) error
}

type TUNDrive struct {
	ctx                context.Context
	lock               sync.Locker
	FD                 uint64
	TunInfo            TUNInfo
	ReadPacketMange    TUNReadPacketManage
	WritePacketManage  TUNWritePacketManage
	ReadPacketChannle  chan []byte
	WritePacketChannle chan []byte
}

func NewDefalutTunDrive(context context.Context) TUNDrive {
	return TUNDrive{
		ctx:                context,
		lock:               &sync.RWMutex{},
		ReadPacketMange:    NewDefalutReadPacketManger(),
		WritePacketManage:  NewDefalutWitePacketManger(),
		ReadPacketChannle:  make(chan []byte, 100000),
		WritePacketChannle: make(chan []byte, 100000),
	}
}

func (tun *TUNDrive) NewTun(t TUNInfo) (err error) {
	if err = t.Check(); err != nil {
		return
	}
	tun.TunInfo = t
	return
}
