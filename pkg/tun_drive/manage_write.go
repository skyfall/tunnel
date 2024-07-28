package tundrive

import "fmt"

type DefalutWritePacketStruct struct {
}

func NewDefalutWitePacketManger() DefalutWritePacketStruct {
	return DefalutWritePacketStruct{}
}

func (r DefalutWritePacketStruct) GetPacket() (ret []byte, err error) {
	var ok bool
	ret, ok = PacketPool.Get().([]byte)
	if ok {
		return
	}
	return nil, fmt.Errorf("get packe err")
}

func (r DefalutWritePacketStruct) DestructionPacket(b []byte) error {
	PacketPool.Put(b)
	return nil
}
