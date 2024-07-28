package tundrive

import "sync"

var PacketPool = &sync.Pool{
	New: func() any {
		return make([]byte, 1500)
	},
}
