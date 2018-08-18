package main

import "crypto/tls"
import "fmt"
import "io"
import "net"
import "os"
import "log"

//import "strings"

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

func main() {
	// Checking if the cert is valid
	cert, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	if err != nil {
		fmt.Println("Something is wrong with Certificate, please check")
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := net.Listen("tcp", ":10053")
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
		conn := tls.Server(connp, config)
		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			fmt.Println()
			c.Close()
		}(conn)
	}
}
