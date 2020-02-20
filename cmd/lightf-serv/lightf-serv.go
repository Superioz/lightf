package main

import (
	"os"

	"github.com/testy/lightf/server"
)

func main() {
	if err := server.Run(); err != nil {
		os.Exit(1)
	}
}
