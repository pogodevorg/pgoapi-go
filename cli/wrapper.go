package cli

import (
	"github.com/urfave/cli"

	"github.com/pkmngo-odi/pogo/api"
	"github.com/pkmngo-odi/pogo/auth"
)

type wrapper struct {
	provider string
	username string
	password string

	lat float64
	lon float64
	alt float64

	debug bool
}

func (w *wrapper) wrap(action func(*cli.Context, *api.Session, auth.Provider) error) func(*cli.Context) error {
	return func(context *cli.Context) error {

		provider, err := auth.NewProvider(w.provider, w.username, w.password)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		location := &api.Location{
			Lon: w.lon,
			Lat: w.lat,
			Alt: w.alt,
		}

		client := api.NewSession(provider, location, w.debug)

		return action(context, client, provider)
	}
}
