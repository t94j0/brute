package brute

type Brute interface {
	Try(host, username, password string) bool
	TryWithPort(host, username, password string, port int) bool
	GetName() string
}

func GetBruteImplMap() map[string]Brute {
	retMap := make(map[string]Brute)

	brutes := []Brute{
		CreateSSHBrute(),
	}

	for _, brute := range brutes {
		retMap[brute.GetName()] = brute
	}

	return retMap
}
