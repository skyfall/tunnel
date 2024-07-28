package tundrive

import (
	"fmt"
)

type DefalutReadPacketStruct struct {
}

func NewDefalutReadPacketManger() DefalutReadPacketStruct {
	return DefalutReadPacketStruct{}
}

func (r DefalutReadPacketStruct) GetPacket() (ret []byte, err error) {
	var ok bool
	ret, ok = PacketPool.Get().([]byte)
	if ok {
		return
	}
	return nil, fmt.Errorf("get packe err")
}

func (r DefalutReadPacketStruct) DestructionPacket(b []byte) error {
	PacketPool.Put(b)
	return nil
}
