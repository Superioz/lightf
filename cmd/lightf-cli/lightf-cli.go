package main

import (
	"os"

	"github.com/testy/lightf/client/cli"
)

func main() {
	if err := cli.Exec(); err != nil {
		os.Exit(1)
	}
}
