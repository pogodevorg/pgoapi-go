package api

import (
	"fmt"
	"log"

	"github.com/pkmngo-odi/pogo/auth"
	"github.com/pkmngo-odi/pogo/rpc"

	ne "github.com/pkmngo-odi/pogo-protos/networking_envelopes"
	nr "github.com/pkmngo-odi/pogo-protos/networking_requests"
	nrm "github.com/pkmngo-odi/pogo-protos/networking_requests_messages"
	nrs "github.com/pkmngo-odi/pogo-protos/networking_responses"

	"github.com/golang/protobuf/proto"
)

const defaultURL = "https://pgorelease.nianticlabs.com/plfe/rpc"

// Session is used to communicate with the Pokémon Go API
type Session struct {
	url      string
	rpc      *rpc.Client
	location *Location
	provider auth.Provider
	debug    bool
}

func generateRequests() []*nr.Request {
	return make([]*nr.Request, 0)
}

// NewSession constructs a Pokémon Go RPC API client
func NewSession(provider auth.Provider, location *Location, debug bool) *Session {
	return &Session{
		rpc:      rpc.NewClient(),
		location: location,
		provider: provider,
		debug:    debug,
	}
}

func (s *Session) setURL(urlToken string) {
	s.url = fmt.Sprintf("https://%s/rpc", urlToken)
}

func (s *Session) getURL() string {
	var url string
	if s.url != "" {
		url = s.url
	} else {
		url = defaultURL
	}
	return url
}

// Call queries the Pokémon Go API through RPC protobuf
func (s *Session) Call(requests []*nr.Request) (*ne.ResponseEnvelope, error) {

	auth := &ne.RequestEnvelope_AuthInfo{
		Provider: s.provider.GetProviderString(),
		Token: &ne.RequestEnvelope_AuthInfo_JWT{
			Contents: s.provider.GetAccessToken(),
			Unknown2: int32(59),
		},
	}

	requestEnvelope := &ne.RequestEnvelope{
		RequestId:  uint64(8145806132888207460),
		StatusCode: int32(2),
		Unknown12:  int64(989),

		Longitude: s.location.Lon,
		Latitude:  s.location.Lat,
		Altitude:  s.location.Alt,

		AuthInfo: auth,

		Requests: requests,
	}

	if s.debug {
		log.Println(proto.MarshalTextString(requestEnvelope))
	}

	responseEnvelope, err := s.rpc.Request(s.getURL(), requestEnvelope)

	if s.debug {
		log.Println(proto.MarshalTextString(responseEnvelope))
	}

	return responseEnvelope, err
}

// Init initializes the client by performing full authentication
func (s *Session) Init() error {
	_, err := s.provider.Login()
	if err != nil {
		return err
	}

	var requests = make([]*nr.Request, 0)
	requests = append(requests, &nr.Request{
		RequestType: nr.RequestType_GET_PLAYER,
	})

	requests = append(requests, &nr.Request{
		RequestType: nr.RequestType_GET_HATCHED_EGGS,
	})

	requests = append(requests, &nr.Request{
		RequestType: nr.RequestType_GET_INVENTORY,
	})

	requests = append(requests, &nr.Request{
		RequestType: nr.RequestType_CHECK_AWARDED_BADGES,
	})

	settingsMessage, _ := proto.Marshal(&nrm.DownloadSettingsMessage{
		Hash: "05daf51635c82611d1aac95c0b051d3ec088a930",
	})

	requests = append(requests, &nr.Request{
		RequestType:    nr.RequestType_DOWNLOAD_SETTINGS,
		RequestMessage: settingsMessage,
	})

	response, err := s.Call(requests)
	if err != nil {
		return err
	}

	url := response.ApiUrl
	if url == "" {
		return fmt.Errorf("Could not initialize session, the service might be down")
	}

	s.setURL(url)
	return nil
}

// GetPlayer returns the current player profile
func (s *Session) GetPlayer() (player *nrs.GetPlayerResponse, err error) {
	requests := generateRequests()
	requests = append(requests, &nr.Request{
		RequestType: nr.RequestType_GET_PLAYER,
	})

	response, err := s.Call(requests)
	if err != nil {
		return player, err
		fmt.Println(response)
	}

	player = &nrs.GetPlayerResponse{}
	proto.Unmarshal(response.Returns[0], player)

	return player, nil
}
