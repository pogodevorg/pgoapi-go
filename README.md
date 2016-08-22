# pogo
Pokémon Go API tools written in Golang

Source at: https://github.com/pokeintel/pogo

## Library
You can include this package as a library in your own project.

### Example
This is an example program that will retrieve player data and print it as JSON.

```go
package main

import (
  "encoding/json"
  "fmt"

  "golang.org/x/net/context"

  "github.com/pokeintel/pogo/api"
  "github.com/pokeintel/pogo/auth"
)

func main() {

  // Unless you already have another net/context complient context, use this empty context.
  // Read more about context at: https://godoc.org/golang.org/x/net/context
  ctx := context.Background()

  // Initialize a new authentication provider to retrieve an access token
  provider, err := auth.NewProvider("ptc", "MyUser", "Pass1!!")
  if err != nil {
    fmt.Println(err)
    return
  }

  // Set the coordinates from where you're connecting
  location := &api.Location{
    Lon: 0.0,
    Lat: 0.0,
    Alt: 0.0,
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
$ go get -u github.com/pokeintel/pogo
```

### Usage

#### Get player profile

```bash
$ pogo -u <username> -p <Secret1234> --lat 0.0 --lon 0.0 player
```

#### Configure through environment variables

```bash
export POGO_ACCOUNT_USERNAME=MyUserAccount
export POGO_ACCOUNT_USERNAME=PasswordThatIsSecret
$ pogo --lat 0.0 --lon 0.0 player
```

## Credit
- Thanks to https://github.com/tejado/pgoapi for inspiration about implementation details.
- Thanks to https://github.com/AeonLucid/POGOProtos for maintaing and constantly improving a Pokémon Go API protobuf specification.
