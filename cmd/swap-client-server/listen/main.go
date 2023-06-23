package main

import (
	"flag"
	"log"
	"net"

	tls "github.com/refraction-networking/utls"
)

func main() {
	var laddrStr string
	flag.StringVar(&laddrStr, "laddr", "0.0.0.0:4443", "local address to connect with")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", laddrStr)
	if err != nil {
		log.Fatal(err)
	}

	listen(laddr)
}

func listen(laddr *net.TCPAddr) error {
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	tcpConn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	conn := tls.UClient(tcpConn, &tls.Config{InsecureSkipVerify: true}, tls.HelloGolang)
	defer conn.Close()

	conn.Write([]byte("Heyyyyyyy"))

	return nil
}
