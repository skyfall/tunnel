package tundrive

import (
	"fmt"
	"net"
	"sync"
)

type TUNPacket struct {
	DesIP    net.IP
	DesPort  int64
	Protocol int
	RawByte  []byte
	EncByte  []byte
	DecByte  []byte
}

var PacketPool = &sync.Pool{
	New: func() any {
		return TUNPacket{
			RawByte:  make([]byte, 0, 1500),
			EncByte:  make([]byte, 0, 1500),
			DecByte:  make([]byte, 0, 1500),
			DesIP:    net.IP{},
			DesPort:  0,
			Protocol: 0,
		}
	},
}

type ManagePakcetStruct struct {
}

func NewDefalutPacketManger() ManagePakcetStruct {
	return ManagePakcetStruct{}
}

func (p ManagePakcetStruct) GetPacket() (ret TUNPacket, err error) {
	var ok bool
	ret, ok = PacketPool.Get().(TUNPacket)
	if ok {
		return
	}

	ret.DecByte = ret.DecByte[:0]
	ret.RawByte = ret.RawByte[:0]
	ret.EncByte = ret.EncByte[:0]
	ret.DesIP = net.IP{}
	ret.DesPort = 0
	ret.Protocol = 0
	err = fmt.Errorf("get packe err")
	return
}

func (r ManagePakcetStruct) DestructionPacket(ret TUNPacket) error {
	PacketPool.Put(ret)
	return nil
}
