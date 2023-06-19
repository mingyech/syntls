package main

import (
	"flag"
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

	c, err := net.ListenPacket("ip4:udp", laddr.IP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	r, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}

	header := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + 20, // 20 bytes for payload
		TTL:      64,
		Protocol: 6, // TCP
		Dst:      raddr.IP,
		Src:      laddr.IP,
		ID:       54321,
	}

	ip := headerToLayer(header)

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

	err = gopacket.SerializeLayers(buf, opts, ip, tcp)
	if err != nil {
		panic(err)
	}

	err = r.WriteTo(header, buf.Bytes(), nil)
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
