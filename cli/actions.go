package cli

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"

	"github.com/pkmngo-odi/pogo/api"
	"github.com/pkmngo-odi/pogo/auth"
)

func fail(e error) *cli.ExitError {
	return cli.NewExitError(e.Error(), 1)
}

func getAccessToken(context *cli.Context, client *api.Session, provider auth.Provider) error {
	token, err := provider.Login()
	if err != nil {
		return fail(err)
	}
	fmt.Println(token)
	return nil
}

func getPlayer(context *cli.Context, client *api.Session, provider auth.Provider) error {
	err := client.Init()
	if err != nil {
		return fail(err)
	}
	profile, err := client.GetPlayer()
	if err != nil {
		return fail(err)
	}
	out, err := json.Marshal(profile)
	if err != nil {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}
