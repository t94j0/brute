package main

import (
	"fmt"

	"github.com/t94j0/brute/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}
}
