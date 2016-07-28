package api

import (
	"fmt"
	"log"
	"time"

	"github.com/pkmngo-odi/pogo/auth"

	"github.com/golang/protobuf/proto"
	"github.com/pkmngo-odi/pogo-protos"
)

const defaultURL = "https://pgorelease.nianticlabs.com/plfe/rpc"
const downloadSettingsHash = "05daf51635c82611d1aac95c0b051d3ec088a930"

// Session is used to communicate with the Pokémon Go API
type Session struct {
	location *Location
	feed     *Feed
	rpc      *RPC
	url      string
	debug    bool

	provider auth.Provider
}

func generateRequests() []*protos.Request {
	return make([]*protos.Request, 0)
}

// NewSession constructs a Pokémon Go RPC API client
func NewSession(provider auth.Provider, location *Location, feed *Feed, debug bool) *Session {
	return &Session{
		location: location,
		rpc:      NewRPC(),
		provider: provider,
		debug:    debug,
		feed:     feed,
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
func (s *Session) Call(requests []*protos.Request) (*protos.ResponseEnvelope, error) {

	auth := &protos.RequestEnvelope_AuthInfo{
		Provider: s.provider.GetProviderString(),
		Token: &protos.RequestEnvelope_AuthInfo_JWT{
			Contents: s.provider.GetAccessToken(),
			Unknown2: int32(59),
		},
	}

	requestEnvelope := &protos.RequestEnvelope{
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

// MoveTo sets your current location
func (s *Session) MoveTo(location *Location) {
	s.location = location
}

// Init initializes the client by performing full authentication
func (s *Session) Init() error {
	_, err := s.provider.Login()
	if err != nil {
		return err
	}

	requests := generateRequests()
	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_PLAYER,
	})

	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_HATCHED_EGGS,
	})

	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_INVENTORY,
	})

	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_CHECK_AWARDED_BADGES,
	})

	settingsMessage, _ := proto.Marshal(&protos.DownloadSettingsMessage{
		Hash: downloadSettingsHash,
	})

	requests = append(requests, &protos.Request{
		RequestType:    protos.RequestType_DOWNLOAD_SETTINGS,
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

// Announce publishes the player's presence and returns the map environment
func (s *Session) Announce() (mapObjects *protos.GetMapObjectsResponse, err error) {

	cellIDs := s.location.GetCellIDs()
	lastTimestamp := time.Now().Unix() * 1000

	requests := generateRequests()

	// Request the map objects based on my current location and route cell ids
	getMapObjectsMessage, _ := proto.Marshal(&protos.GetMapObjectsMessage{
		// Traversed route since last supposed last heartbeat
		CellId: cellIDs,

		// Timestamps in milliseconds corresponding to each route cell id
		SinceTimestampMs: make([]int64, len(cellIDs)),

		// Current longitide and latitude
		Longitude: s.location.Lon,
		Latitude:  s.location.Lat,
	})

	requests = append(requests, &protos.Request{
		RequestType:    protos.RequestType_GET_MAP_OBJECTS,
		RequestMessage: getMapObjectsMessage,
	})

	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_HATCHED_EGGS,
	})

	// Request the inventory with a message containing the current time
	getInventoryMessage, _ := proto.Marshal(&protos.GetInventoryMessage{
		LastTimestampMs: lastTimestamp,
	})

	requests = append(requests, &protos.Request{
		RequestType:    protos.RequestType_GET_INVENTORY,
		RequestMessage: getInventoryMessage,
	})

	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_CHECK_AWARDED_BADGES,
	})

	settingsMessage, _ := proto.Marshal(&protos.DownloadSettingsMessage{
		Hash: downloadSettingsHash,
	})

	requests = append(requests, &protos.Request{
		RequestType:    protos.RequestType_DOWNLOAD_SETTINGS,
		RequestMessage: settingsMessage,
	})

	response, err := s.Call(requests)
	if err != nil {
		return mapObjects, &RequestError{}
	}

	mapObjects = &protos.GetMapObjectsResponse{}
	proto.Unmarshal(response.Returns[0], mapObjects)
	s.feed.Push(mapObjects)

	return mapObjects, GetErrorFromStatus(response.StatusCode)
}

// GetPlayer returns the current player profile
func (s *Session) GetPlayer() (player *protos.GetPlayerResponse, err error) {
	requests := generateRequests()
	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_PLAYER,
	})

	response, err := s.Call(requests)
	if err != nil {
		return player, err
	}

	player = &protos.GetPlayerResponse{}
	proto.Unmarshal(response.Returns[0], player)
	s.feed.Push(player)

	return player, nil
}

// GetInventory returns the player items
func (s *Session) GetInventory() (inventory *protos.GetInventoryResponse, err error) {
	requests := generateRequests()
	requests = append(requests, &protos.Request{
		RequestType: protos.RequestType_GET_INVENTORY,
	})

	response, err := s.Call(requests)
	if err != nil {
		return inventory, err
	}

	inventory = &protos.GetInventoryResponse{}
	proto.Unmarshal(response.Returns[0], inventory)
	s.feed.Push(inventory)

	return inventory, nil
}
