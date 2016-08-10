package main

import (
	"os"

	"github.com/pokeintel/pogo/api"
	"github.com/pokeintel/pogo/cli"
)

func main() {
	crypto := &api.DefaultCrypto{}
	cli.Run(crypto, os.Args)
}
