package brute

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/dutchcoders/goftp"
)

type FTPBrute struct {
	Protocol string `default:"ssh"`
	Port     string `cli:"port" default:"22" required:"true"`
}

func (s FTPBrute) Try(host, username string, password []byte) bool {
	host += ":" + s.Port

	ftp, err := goftp.Connect(host)
	if err != nil {
		fmt.Println("Error connecting to host")
		os.Exit(1)
	}

	config := tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequestClientCert,
	}

	//
	if err := ftp.AuthTLS(&config); err != nil {
		fmt.Println("Error: Could not get TLS")
	}

	// Username / password authentication
	return ftp.Login(username, string(password)) == nil

	return err == nil
}

func (s FTPBrute) GetProtocol() string {
	return s.Protocol
}
