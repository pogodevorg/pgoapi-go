package main

import (
	"fmt"
	"os"

	"github.com/zeeraw/pogo/api"

	"github.com/urfave/cli"
)

var client *api.Client
var provider string
var username string
var password string

func getClient() (*api.Client, error) {
	if client != nil {
		return client, nil
	}

	var creds api.Credentials
	switch provider {
	case "ptc":
		creds = &api.PTCCredentials{
			Username: username,
			Password: password,
		}
	default:
		return nil, cli.NewExitError(fmt.Sprintf("Provider \"%s\" is not supported", provider), 1)
	}

	client = api.NewClient(creds)
	return client, nil
}

func getAccessToken(context *cli.Context, client *api.Client) error {
	token, err := client.GetAccessToken()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Println(token)
	return nil
}

func wrap(action func(*cli.Context, *api.Client) error) func(*cli.Context) error {
	return func(context *cli.Context) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		return action(context, client)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "pogo"
	app.Usage = "Command line client for the Pokémon Go API"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username,u",
			Destination: &username,
			EnvVar:      "POGO_ACCOUNT_USERNAME",
		},
		cli.StringFlag{
			Name:        "password,p",
			Destination: &password,
			EnvVar:      "POGO_ACCOUNT_PASSWORD",
		},
		cli.StringFlag{
			Name:        "provider",
			Destination: &provider,
			Value:       "ptc",
			Usage:       "Your account provider can be either \"ptc\" or \"google\"",
			EnvVar:      "POGO_ACCOUNT_PROVIDER",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "get_access_token",
			Usage:  "Retrieves an API access token from your credentials",
			Action: wrap(getAccessToken),
		},
	}
	app.Run(os.Args)
}
