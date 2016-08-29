package cli

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"

	"github.com/urfave/cli"

	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/pogodevorg/pgoapi-go/auth"
)

func fail(e error) *cli.ExitError {
	return cli.NewExitError(e.Error(), 1)
}

func isFailure(e error) bool {
	if e != nil {
		switch e {
		case api.ErrNewRPCURL:
			return false
		default:
			return true
		}
	}
	return false
}

func getAccessToken(ctx context.Context, session *api.Session, provider auth.Provider) error {
	token, err := provider.Login(ctx)
	if isFailure(err) {
		return fail(err)
	}
	fmt.Println(token)
	return nil
}

func getPlayer(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if isFailure(err) {
		return fail(err)
	}
	profile, err := session.GetPlayer(ctx)
	if isFailure(err) {
		return fail(err)
	}
	out, err := json.Marshal(profile)
	if isFailure(err) {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}

func getInventory(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if isFailure(err) {
		return fail(err)
	}
	inventory, err := session.GetInventory(ctx)
	if isFailure(err) {
		return fail(err)
	}
	out, err := json.Marshal(inventory)
	if isFailure(err) {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}

func getMap(ctx context.Context, session *api.Session, provider auth.Provider) error {
	err := session.Init(ctx)
	if isFailure(err) {
		return fail(err)
	}
	mapObjects, err := session.GetPlayerMap(ctx)
	if isFailure(err) {
		return fail(err)
	}
	out, err := json.Marshal(mapObjects)
	if isFailure(err) {
		return fail(err)
	}

	fmt.Println(string(out))
	return nil
}
