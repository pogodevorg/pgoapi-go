package cli

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"

	"github.com/urfave/cli"

	"github.com/pkmngo-odi/pogo/api"
	"github.com/pkmngo-odi/pogo/auth"
)

func fail(e error) *cli.ExitError {
	return cli.NewExitError(e.Error(), 1)
}

func getAccessToken(ctx context.Context, session *api.Session, provider auth.Provider) error {
	token, err := provider.Login(ctx)
	if err != nil {
		return fail(err)
	}
	fmt.Println(token)
	return nil
}

func getPlayer(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if err != nil {
		return fail(err)
	}
	profile, err := session.GetPlayer(ctx)
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

func getInventory(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if err != nil {
		return fail(err)
	}
	inventory, err := session.GetInventory(ctx)
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

func getMap(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if err != nil {
		return fail(err)
	}
	mapObjects, err := session.GetPlayerMap(ctx)
	if err != nil {
		return fail(err)
	}
	out, err := json.Marshal(mapObjects)
	if err != nil {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}
