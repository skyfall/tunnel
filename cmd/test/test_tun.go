package main

import (
	"net"
	tundrive "tunnel/pkg/tun_drive"
)

func main() {

	tuninfo := tundrive.TUNInfo{}
	tuninfo.Name = "tun_1"
	tuninfo.LocalIP = net.ParseIP("192.168.100.22")
	tuninfo.Mask = net.ParseIP("255.255.0.0")
	tuninfo.Mtu = 1111
	tuninfo.Type = tundrive.TUN_TYPE_TUN

	err := tuninfo.CreateTunDrive()
	if err != nil {
		println(err.Error())
	}

	tuninfo.RunTunDrive()

	println("fd :")
}
