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

	ip := &layers.IPv4{
		SrcIP:    net.ParseIP(c.LocalAddr().String()),
		DstIP:    net.ParseIP(c.RemoteAddr().String()),
		Protocol: layers.IPProtocolTCP,
	}

	fmt.Printf("ip header: %+v\n", ip)

	// Define the TCP layer
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(laddr.Port),
		DstPort: layers.TCPPort(raddr.Port),
		SYN:     true,
		Seq:     11050,
		Window:  14600,
	}
	tcp.SetNetworkLayerForChecksum(ip)

	// Stack all layers and serialize them
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	err = gopacket.SerializeLayers(buf, opts, tcp)
	if err != nil {
		panic(err)
	}

	// err = r.WriteTo(header, buf.Bytes(), nil)
	c.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
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
