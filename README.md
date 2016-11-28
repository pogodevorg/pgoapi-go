# pgoapi-go
Pokémon Go API tools written in Golang

[![Build Status](https://travis-ci.org/pogodevorg/pgoapi-go.svg?branch=master)](https://travis-ci.org/pogodevorg/pgoapi-go)

Source at: https://github.com/pogodevorg/pgoapi-go

## Dependencies
Only supports Go `1.7`

## Library
You can include this package as a library in your own project.

### Example
This is an example program that will retrieve player data and print it as JSON.

```go
package main

import (
  "encoding/json"
  "fmt"
  "context"

  "github.com/pogodevorg/pgoapi-go/api"
  "github.com/pogodevorg/pgoapi-go/auth"
)

func main() {

  // Unless you already have another net/context compliant context, use this empty context.
  // Read more about context at: https://godoc.org/context
  ctx := context.Background()

  // Initialize a new authentication provider to retrieve an access token
  provider, err := auth.NewProvider("ptc", "MyUser", "Pass1!!")
  if err != nil {
    fmt.Println(err)
    return
  }

  // If you have previously retrieved a token, you can also set token.

  provider.SetAccessToken("cafedecafbeefface");

  // Set the coordinates from where you're connecting
  location := &api.Location{
    Lon: 0.0,
    Lat: 0.0,
    Alt: 0.0,
    Accuracy: 3.0,
  }

  // Set up a feed to where all the responses will be pushed
  // The void feed will do nothing with the response entries
  feed := &api.VoidFeed{}

  // Set up the type of crypto to use for signing requests
  //
  // For most intents and purposes, you should be fine with
  // using the Default crypto.
  crypto := &api.DefaultCrypto{}

  // Start new session and connect
  session := api.NewSession(provider, location, feed, crypto, false)

  // (If you have previously set a token, Init will not log you in but instead do request using that token.)
  err = session.Init(ctx)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Start querying the API
  player, err := session.GetPlayer(ctx)
  if err != nil {
    fmt.Println(err)
    return
  }

  out, err := json.Marshal(player)
  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println(string(out))
}
```

### Using the feed
The feed is a common interface to get a stream of all responses.
This debug feed will print all wild pokemon and forts from map responses to standard out.

```go
type DebugFeed struct {}

func (f *DebugFeed) Push(entry interface{}) {
  switch e := entry.(type) {
  default:
    // NOOP: Will not report type
  case *protos.GetMapObjectsResponse:
    cells := e.GetMapCells()
    for _, cell := range cells {
      pokemons := cell.GetWildPokemons()
      if len(pokemons) > 0 {
        fmt.Println(pokemons)
      }
      forts := cell.GetForts()
      if len(forts) > 0 {
        fmt.Println(forts)
      }
    }
  }
}
```

## Command line tool

### Install
Make sure you're running the latest version of Go, then simply install through `go get`.

```bash
$ go get -u github.com/pogodevorg/pgoapi-go
```

### Usage

#### Get player profile

```bash
$ pgoapi-go -u <username> -p <Secret1234> --lat 0.0 --lon 0.0 player
```

#### Configure through environment variables

```bash
export PGOAPI_ACCOUNT_USERNAME=MyUserAccount
export PGOAPI_ACCOUNT_PASSWORD=PasswordThatIsSecret
$ pgoapi-go --lat 0.0 --lon 0.0 player
```

## Credit
- Thanks to https://github.com/tejado/pgoapi for inspiration about implementation details.
- Thanks to https://github.com/AeonLucid/POGOProtos for maintaining and constantly improving a Pokémon Go API Protobuf specification.
