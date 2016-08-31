package cli

import (
	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/urfave/cli"
)

// Run interprets arguments and performs actions
func Run(crypto api.Crypto, args []string) {

	w := wrapper{
		crypto: crypto,
	}

	app := cli.NewApp()
	app.Name = "pgoapi-go"
	app.Usage = "Command line client for the Pokémon Go API"
	app.Author = "Philip Vieira <zee@vall.in>"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug,d",
			Destination: &w.debug,
			EnvVar:      "PGOAPI_DEBUG",
		},
		cli.StringFlag{
			Name:        "username,u",
			Destination: &w.username,
			EnvVar:      "PGOAPI_ACCOUNT_USERNAME",
		},
		cli.StringFlag{
			Name:        "password,p",
			Destination: &w.password,
			EnvVar:      "PGOAPI_ACCOUNT_PASSWORD",
		},
		cli.StringFlag{
			Name:        "provider",
			Destination: &w.provider,
			Value:       "ptc",
			Usage:       "Your account provider can be either \"ptc\" or \"google\"",
			EnvVar:      "PGOAPI_ACCOUNT_PROVIDER",
		},
		cli.Float64Flag{
			Name:        "latitude,lat",
			Destination: &w.lat,
			Value:       0.0,
			EnvVar:      "PGOAPI_DEFAULT_LATITUDE",
		},
		cli.Float64Flag{
			Name:        "longitude,lon",
			Destination: &w.lon,
			Value:       0.0,
			EnvVar:      "PGOAPI_DEFAULT_LONGITUDE",
		},
		cli.Float64Flag{
			Name:        "altitude,alt",
			Destination: &w.alt,
			Value:       0.0,
			EnvVar:      "PGOAPI_DEFAULT_ALTITUDE",
		},
		cli.Float64Flag{
			Name:        "accuracy,acc",
			Destination: &w.accuracy,
			Value:       3.0,
			EnvVar:      "PGOAPI_DEFAULT_ACCURACY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "access_token",
			Usage:  "Retrieves an API access token from your credentials",
			Action: w.wrap(getAccessToken),
		},
		{
			Name:   "player",
			Usage:  "Retrieves the user's Pokémon Go player profile",
			Action: w.wrap(getPlayer),
		},
		{
			Name:   "inventory",
			Usage:  "Retrieves the user's Pokémon Go player inventory",
			Action: w.wrap(getInventory),
		},
		{
			Name:   "map",
			Usage:  "Retrieves map data for the player's current location",
			Action: w.wrap(getMap),
		},
	}

	app.Run(args)
}
