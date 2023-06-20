package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

func main() {

	var raddrStr string
	var laddrStr string

	flag.StringVar(&raddrStr, "raddr", "127.0.0.1:443", "remote address to connect to")
	flag.StringVar(&laddrStr, "laddr", "0.0.0.0:443", "local address to connect with")
	flag.Parse()

	raddr, err := net.ResolveTCPAddr("tcp", raddrStr)
	if err != nil {
		log.Fatal(err)
	}

	laddr, err := net.ResolveTCPAddr("tcp", laddrStr)
	if err != nil {
		log.Fatal(err)
	}

	c, err := net.DialIP("ip4:tcp", &net.IPAddr{IP: laddr.IP}, &net.IPAddr{IP: raddr.IP})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	buf := make([]byte, 256)
	_, err = c.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	packet := gopacket.NewPacket(buf, layers.LayerTypeTCP, gopacket.Default)
	tcpLayer := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	fmt.Printf("received: %+v\n", tcpLayer)

}

func headerToLayer(header *ipv4.Header) *layers.IPv4 {
	return &layers.IPv4{
		Version:    4,
		IHL:        uint8(header.Len / 4), // header length is in 4-octet units
		TOS:        uint8(header.TOS),     // Type of Service
		Length:     uint16(header.TotalLen),
		Id:         uint16(header.ID),
		Flags:      layers.IPv4DontFragment, // Depends on your header settings
		FragOffset: uint16(header.FragOff),
		TTL:        uint8(header.TTL),
		Protocol:   layers.IPProtocol(header.Protocol),
		SrcIP:      header.Src,
		DstIP:      header.Dst,
	}
}
