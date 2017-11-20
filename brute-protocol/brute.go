package brute

import (
	"fmt"
	"reflect"
)

type Brute interface {
	Try(host, username string, password []byte) bool
	GetProtocol() string
}

func GetBruteMap() map[string]Brute {
	bruteMap := make(map[string]Brute)

	list := []Brute{
		SSHBrute{},
	}

	// Get default values for protocols by reflection
	for _, bruteStruct := range list {
		bruteType := reflect.TypeOf(bruteStruct)
		for i := 0; i < bruteType.NumField(); i++ {
			bf := bruteType.Field(i)
			if bf.Name == "Protocol" {
				protocolName := bf.Tag.Get("default")
				bruteMap[protocolName] = bruteStruct
			}
		}
	}

	return bruteMap
}

func GetCliHelp(module Brute) string {
	output := ""

	moduleType := reflect.TypeOf(module)
	for i := 0; i < moduleType.NumField(); i++ {
		if moduleType.Field(i).Name == "Protocol" {
			continue
		}

		tag := moduleType.Field(i).Tag
		cli := tag.Get("cli")
		def, hasDef := tag.Lookup("default")
		req, isReq := tag.Lookup("required")

		output += fmt.Sprintf("\t--%s ", cli)

		if hasDef {
			output += fmt.Sprintf("(default: %s) ", def)
		}
		if isReq && req == "true" {
			output += fmt.Sprintf("REQUIRED")
		}
		output += "\n"

	}

	return output
}
