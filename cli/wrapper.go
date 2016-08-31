package cli

import (
	"context"

	"github.com/urfave/cli"

	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/pogodevorg/pgoapi-go/auth"
)

type wrapper struct {
	provider string
	username string
	password string

	lat      float64
	lon      float64
	alt      float64
	accuracy float64

	debug  bool
	crypto api.Crypto
}

func (w *wrapper) wrap(action func(context.Context, *api.Session, auth.Provider) error) func(*cli.Context) error {
	return func(c *cli.Context) error {

		ctx := context.Background()

		provider, err := auth.NewProvider(w.provider, w.username, w.password)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		location := &api.Location{
			Lon:      w.lon,
			Lat:      w.lat,
			Alt:      w.alt,
			Accuracy: w.accuracy,
		}

		client := api.NewSession(provider, location, &api.VoidFeed{}, w.crypto, w.debug)

		return action(ctx, client, provider)
	}
}
