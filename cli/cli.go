package cli

import (
	"github.com/urfave/cli"
)

// Run interprets arguments and performs actions
func Run(args []string) {

	w := wrapper{}

	app := cli.NewApp()
	app.Name = "pogo"
	app.Usage = "Command line client for the Pokémon Go API"
	app.Author = "Philip Vieira <zee@vall.in>"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug,d",
			Destination: &w.debug,
			EnvVar:      "POGO_DEBUG",
		},
		cli.StringFlag{
			Name:        "username,u",
			Destination: &w.username,
			EnvVar:      "POGO_ACCOUNT_USERNAME",
		},
		cli.StringFlag{
			Name:        "password,p",
			Destination: &w.password,
			EnvVar:      "POGO_ACCOUNT_PASSWORD",
		},
		cli.StringFlag{
			Name:        "provider",
			Destination: &w.provider,
			Value:       "ptc",
			Usage:       "Your account provider can be either \"ptc\" or \"google\"",
			EnvVar:      "POGO_ACCOUNT_PROVIDER",
		},
		cli.Float64Flag{
			Name:        "latitude,lat",
			Destination: &w.lat,
			Value:       0.0,
			EnvVar:      "POGO_DEFAULT_LATITUDE",
		},
		cli.Float64Flag{
			Name:        "longitude,lon",
			Destination: &w.lon,
			Value:       0.0,
			EnvVar:      "POGO_DEFAULT_LONGITUDE",
		},
		cli.Float64Flag{
			Name:        "altitude,alt",
			Destination: &w.alt,
			Value:       0.0,
			EnvVar:      "POGO_DEFAULT_ALTITUDE",
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
	}

	app.Run(args)
}
