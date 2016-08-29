package main

import (
	"os"

	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/pogodevorg/pgoapi-go/cli"
)

func main() {
	crypto := &api.DefaultCrypto{}
	cli.Run(crypto, os.Args)
}
