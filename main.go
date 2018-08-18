package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

const serverKey = `-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIFbNiJ54CxDrAPZ1SMNoxzqA0kjZOvfrlqzkWJFijqroAoGCCqGSM49
AwEHoUQDQgAEPrHq9gV0RsOEaPz7LxWrDkf7O7r7geyHImHfLK6iYsXNKTtDJayd
ReTmha0G+V5a3fKiJfxQgU69hyPgMW60tA==
-----END EC PRIVATE KEY-----
`

const serverCert = `-----BEGIN CERTIFICATE-----
MIIClzCCAj2gAwIBAgIJALbSpDn+Z1c9MAkGByqGSM49BAEwaTELMAkGA1UEBhMC
REUxDzANBgNVBAgTBkJlcmxpbjEPMA0GA1UEBxMGQmVybGluMREwDwYDVQQKEwhQ
aGlsIENvLjElMCMGCSqGSIb3DQEJARYWcGNhcnZhbGhvLnRpQGdtYWlsLmNvbTAe
Fw0xODA4MTgxMzE1MTFaFw0yODA4MTUxMzE1MTFaMGkxCzAJBgNVBAYTAkRFMQ8w
DQYDVQQIEwZCZXJsaW4xDzANBgNVBAcTBkJlcmxpbjERMA8GA1UEChMIUGhpbCBD
by4xJTAjBgkqhkiG9w0BCQEWFnBjYXJ2YWxoby50aUBnbWFpbC5jb20wWTATBgcq
hkjOPQIBBggqhkjOPQMBBwNCAAQ+ser2BXRGw4Ro/PsvFasOR/s7uvuB7IciYd8s
rqJixc0pO0MlrJ1F5OaFrQb5Xlrd8qIl/FCBTr2HI+AxbrS0o4HOMIHLMB0GA1Ud
DgQWBBRCtzkKZuixIMZjXPR6zQOxgsNmYTCBmwYDVR0jBIGTMIGQgBRCtzkKZuix
IMZjXPR6zQOxgsNmYaFtpGswaTELMAkGA1UEBhMCREUxDzANBgNVBAgTBkJlcmxp
bjEPMA0GA1UEBxMGQmVybGluMREwDwYDVQQKEwhQaGlsIENvLjElMCMGCSqGSIb3
DQEJARYWcGNhcnZhbGhvLnRpQGdtYWlsLmNvbYIJALbSpDn+Z1c9MAwGA1UdEwQF
MAMBAf8wCQYHKoZIzj0EAQNJADBGAiEA5hysuLyB6JvwsYn4KpdJySvkRg4lc61Q
BbXaxBR98RACIQDNvEj2uojQfY2sV94kkJJxAAMXwejQpji6wjEEkOSqDg==
-----END CERTIFICATE-----
`

const loAddr string = "127.0.0.1:10053"
const reAddr string = "cloudflare-dns.com:853"

func main() {
	// Checking if the cert is valid
	cert, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	if err != nil {
		fmt.Println("Something is wrong with Certificate, please check")
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	//Opening Listener
	ln, err := tls.Listen("tcp", loAddr, config)
	if err != nil {
		fmt.Println("Listeners Could not be opened")
	}
	defer ln.Close()

	// Giving some feedback
	fmt.Println("Launching Local DNS server...")

	// accept connection on port
	fmt.Println("DNS server is accepting connections")

	// run loop forever (or until ctrl-c)
	for {
		connp, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go proxyFunc(connp)
	}
}

func proxyFunc(conn net.Conn) {
	defer conn.Close()

	rAddr, err := net.ResolveTCPAddr("tcp", reAddr)
	if err != nil {
		fmt.Println("Cannot open connection to remote address")
	}

	reConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		log.Print(err)
	}
	defer reConn.Close()

	Pipeline(conn, reConn)

	fmt.Println("Connection ended with remote address")
}

func channCoon(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

func Pipeline(conn0 net.Conn, conn1 net.Conn) {
	chan0 := channCoon(conn0)
	chan1 := channCoon(conn1)

	for {
		select {
		case b0 := <-chan0:
			if b0 == nil {
				return
			} else {
				conn1.Write(b0)
			}
		case b1 := <-chan1:
			if b1 == nil {
				return
			} else {
				conn0.Write(b1)
			}
		}
	}
}
