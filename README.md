# pogo
Pokémon Go API tools written in Golang
Source at: https://github.com/pkmngo-odi/pogo

## Library
You can include this package as a library in your own project.

### Example

```go
import (
  "fmt"
  "github.com/pkmngo-odi/pogo/api"
  "github.com/pkmngo-odi/pogo/auth"
)

func main() {

  // Initialize a new authentication provider to retrieve an access token
  provider, err := auth.NewProvider("ptc", "MyUser", "Pass1!!)")
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

  // Start new session and connect
  session := api.NewSession(provider, location, false)
  session.Init()

  // Start querying the API
  player, err := client.GetPlayer()
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
  return nil

}
```

## Command line tool

### Install
Make sure you're running the latest version of Go, then simply install through `go get`.

```bash
$ go get -u github.com/pkmngo-odi/pogo
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
