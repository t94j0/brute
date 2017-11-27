package cli

import (
	"bufio"
	"os"
)

type WordListBrute struct {
	Password        string
	PasswordList    *bufio.Reader
	PasswordInitial *os.File
	PasswordDone    bool

	Username        string
	UsernameList    *bufio.Reader
	UsernameInitial *os.File
	UsernameDone    bool

	Host        string
	HostList    *bufio.Reader
	HostInitial *os.File
}

func openList(path string) (*bufio.Reader, *os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewReader(file), file, nil
}

// NewWordListBrute fills in the WordListBrute struct with opened files and
// bufio readers to make the next password selection process easier
//
// h is the host, hp is the host path, u is the username, up is the username
// path p is the password, pp is the password path
func NewWordListBrute(h, hp, u, up, p, pp string) (wl *WordListBrute, err error) {
	wl = new(WordListBrute)

	if hp == "" {
		wl.Host = h
	} else {
		wl.HostList, wl.HostInitial, err = openList(hp)
		if err != nil {
			return nil, err
		}
	}

	if up == "" {
		wl.Username = u
	} else {
		wl.UserList, wl.UsernameInitial, err = openList(up)
		if err != nil {
			return nil, err
		}
	}

	if pp == "" {
		wl.Password = p
	} else {
		wl.PasswordList, wl.PasswordInitial, err = openList(pp)
		if err != nil {
			return nil, err
		}
	}

	return wl, nil
}

// nextPassword is the next password in the list. The first returned value
// is if the wordlist has another password or not. The second value is if
// the next value has a password
func (w *WordListBrute) NextUsername() (bool, string) {
	if w.UsernameDone {
		return false, ""
	} else if w.UsernameList == nil {
		w.UsernameDone = true
		return w.Username, true
	} else {
		nextString, err := w.UsernameList.ReadString('\n')
		if err != nil {
			return false, ""
		}
		return true, nextString
	}
}

func (w *WordListBrute) NextPassword() (bool, string) {
	if w.PasswordDone {
		return false, ""
	} else if w.PasswordList == nil {
		w.PasswordDone = true
		return w.Password, true
	} else {
		nextString, err := w.PasswordList.ReadString('\n')
		if err != nil {
			return false, ""
		}
		return true, nextString
	}
}

func (w *WordListBrute) NextHost() string {
	if w.HostList == nil {
		tmpHost := w.Host
		w.Host = ""
		return tmpHost
	} else {
		nextString, err := w.HostList.ReadString('\n')
		if err != nil {
			return ""
		}
		return nextString
	}
}
