package main

import (
	"fmt"
	"net"
)

func main() {
	mac_address := "0C:8B:7D:25:38:35"
	mac_bytes , _ := net.ParseMAC( mac_address )
	magic_packet := []byte{}
	for i := 0; i < 6; i++ {
		magic_packet = append( magic_packet , 0xFF )
	}
	for i := 0; i < 16; i++ {
		magic_packet = append( magic_packet , mac_bytes... )
	}
	addr := &net.UDPAddr{
		IP:   net.IPv4bcast ,
		Port: 9 ,
	}
	conn , _ := net.DialUDP( "udp" , nil , addr )
	defer conn.Close()
	result , err := conn.Write( magic_packet )
	fmt.Println( result , err )
}