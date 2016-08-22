package main

import (
	"os"

	"github.com/pogodevorg/pogo/api"
	"github.com/pogodevorg/pogo/cli"
)

func main() {
	crypto := &api.DefaultCrypto{}
	cli.Run(crypto, os.Args)
}
