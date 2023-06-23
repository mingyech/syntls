package main

import (
	"crypto/x509"
	"encoding/hex"
	"flag"
	"log"
	"net"

	tls "github.com/refraction-networking/utls"
)

func fromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

var testP256PrivateKey, _ = x509.ParseECPrivateKey(fromHex("30770201010420012f3b52bc54c36ba3577ad45034e2e8efe1e6999851284cb848725cfe029991a00a06082a8648ce3d030107a14403420004c02c61c9b16283bbcc14956d886d79b358aa614596975f78cece787146abf74c2d5dc578c0992b4f3c631373479ebf3892efe53d21c4f4f1cc9a11c3536b7f75"))
var testP256Certificate = fromHex("308201693082010ea00302010202105012dc24e1124ade4f3e153326ff27bf300a06082a8648ce3d04030230123110300e060355040a130741636d6520436f301e170d3137303533313232343934375a170d3138303533313232343934375a30123110300e060355040a130741636d6520436f3059301306072a8648ce3d020106082a8648ce3d03010703420004c02c61c9b16283bbcc14956d886d79b358aa614596975f78cece787146abf74c2d5dc578c0992b4f3c631373479ebf3892efe53d21c4f4f1cc9a11c3536b7f75a3463044300e0603551d0f0101ff0404030205a030130603551d25040c300a06082b06010505070301300c0603551d130101ff04023000300f0603551d1104083006820474657374300a06082a8648ce3d0403020349003046022100963712d6226c7b2bef41512d47e1434131aaca3ba585d666c924df71ac0448b3022100f4d05c725064741aef125f243cdbccaa2a5d485927831f221c43023bd5ae471a")

func main() {
	var raddrStr string
	flag.StringVar(&raddrStr, "raddr", "0.0.0.0:4443", "remote address to connect to")
	flag.Parse()

	raddr, err := net.ResolveTCPAddr("tcp", raddrStr)
	if err != nil {
		log.Fatal(err)
	}

	dial(raddr)
}

func dial(raddr *net.TCPAddr) error {
	tcpConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		panic(err)
	}

	conn := tls.Server(tcpConn,
		&tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{testP256Certificate},
			PrivateKey: testP256PrivateKey,
		}}})
	defer conn.Close()

	buffer := make([]byte, 4096)

	bytesRead, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	} else {
		log.Println("Received message:", string(buffer[:bytesRead]))
	}
	return nil
}
