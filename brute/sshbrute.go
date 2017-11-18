package brute

import (
	"strconv"

	"golang.org/x/crypto/ssh"
)

type SSHBrute struct {
	Name string
}

func CreateSSHBrute() SSHBrute {
	return SSHBrute{
		Name: "ssh",
	}
}

func (s SSHBrute) Try(host, username, password string) bool {
	return s.TryWithPort(host, username, password, 22)
}

func (s SSHBrute) TryWithPort(host, username, password string, port int) bool {
	host += ":" + strconv.Itoa(port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	_, err := ssh.Dial("tcp", host, config)

	return err == nil
}

func (s SSHBrute) GetName() string {
	return s.Name
}
