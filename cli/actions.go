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

func getAccessToken(context *cli.Context, session *api.Session, provider auth.Provider) error {
	token, err := provider.Login()
	if err != nil {
		return fail(err)
	}
	fmt.Println(token)
	return nil
}

func getPlayer(context *cli.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init()
	if err != nil {
		return fail(err)
	}
	profile, err := session.GetPlayer()
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

func getInventory(context *cli.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init()
	if err != nil {
		return fail(err)
	}
	inventory, err := session.GetInventory()
	if err != nil {
		return fail(err)
	}
	out, err := json.Marshal(inventory)
	if err != nil {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}
