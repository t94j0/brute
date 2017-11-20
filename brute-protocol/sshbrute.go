package brute

import (
	"golang.org/x/crypto/ssh"
)

type SSHBrute struct {
	Protocol string `default:"ssh"`
	Port     string `cli:"port" default:"22" required:"true"`
}

func (s SSHBrute) Try(host, username string, password []byte) bool {
	host += ":" + s.Port

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	_, err := ssh.Dial("tcp", host, config)

	return err == nil
}

func (s SSHBrute) GetProtocol() string {
	return s.Protocol
}
