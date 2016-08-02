package api

import (
	"fmt"
	"log"
	"time"

	"github.com/pkmngo-odi/pogo/auth"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkmngo-odi/pogo-protos"
)

const defaultURL = "https://pgorelease.nianticlabs.com/plfe/rpc"
const downloadSettingsHash = "05daf51635c82611d1aac95c0b051d3ec088a930"

// Session is used to communicate with the Pokémon Go API
type Session struct {
	feed     Feed
	location *Location
	rpc      *RPC
	url      string
	debug    bool
	debugger *jsonpb.Marshaler

	provider auth.Provider
}

func generateRequests() []*protos.Request {
	return make([]*protos.Request, 0)
}

// NewSession constructs a Pokémon Go RPC API client
func NewSession(provider auth.Provider, location *Location, feed Feed, debug bool) *Session {
	return &Session{
		location: location,
		rpc:      NewRPC(),
		provider: provider,
		debug:    debug,
		debugger: &jsonpb.Marshaler{Indent: "\t"},
		feed:     feed,
	}
}

// SetTimeout sets the client timeout for the RPC API
func (s *Session) SetTimeout(d time.Duration) {
	s.rpc.http.Timeout = d
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
		log.Println(s.debugger.MarshalToString(requestEnvelope))
	}

	responseEnvelope, err := s.rpc.Request(s.getURL(), requestEnvelope)

	if s.debug {
		log.Println(s.debugger.MarshalToString(responseEnvelope))
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

	settingsMessage, _ := proto.Marshal(&protos.DownloadSettingsMessage{
		Hash: downloadSettingsHash,
	})
	requests := []*protos.Request{
		{RequestType: protos.RequestType_GET_PLAYER},
		{RequestType: protos.RequestType_GET_HATCHED_EGGS},
		{RequestType: protos.RequestType_GET_INVENTORY},
		{RequestType: protos.RequestType_CHECK_AWARDED_BADGES},
		{protos.RequestType_DOWNLOAD_SETTINGS, settingsMessage},
	}

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

	settingsMessage, _ := proto.Marshal(&protos.DownloadSettingsMessage{
		Hash: downloadSettingsHash,
	})
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
	// Request the inventory with a message containing the current time
	getInventoryMessage, _ := proto.Marshal(&protos.GetInventoryMessage{
		LastTimestampMs: lastTimestamp,
	})
	requests := []*protos.Request{
		{RequestType: protos.RequestType_GET_PLAYER},
		{RequestType: protos.RequestType_GET_HATCHED_EGGS},
		{protos.RequestType_GET_INVENTORY, getInventoryMessage},
		{RequestType: protos.RequestType_CHECK_AWARDED_BADGES},
		{protos.RequestType_DOWNLOAD_SETTINGS, settingsMessage},
		{protos.RequestType_GET_MAP_OBJECTS, getMapObjectsMessage},
	}

	response, err := s.Call(requests)
	if err != nil {
		return mapObjects, &RequestError{}
	}

	mapObjects = &protos.GetMapObjectsResponse{}
	err = proto.Unmarshal(response.Returns[0], mapObjects)
	if err != nil {
		return nil, &ResponseError{err}
	}
	s.feed.Push(mapObjects)

	return mapObjects, GetErrorFromStatus(response.StatusCode)
}

// GetPlayer returns the current player profile
func (s *Session) GetPlayer() (*protos.GetPlayerResponse, error) {
	requests := []*protos.Request{{RequestType: protos.RequestType_GET_PLAYER}}
	response, err := s.Call(requests)
	if err != nil {
		return nil, err
	}

	player := &protos.GetPlayerResponse{}
	err = proto.Unmarshal(response.Returns[0], player)
	if err != nil {
		return nil, &ResponseError{err}
	}
	s.feed.Push(player)

	return player, GetErrorFromStatus(response.StatusCode)
}

// GetPlayerMap returns the surrounding map cells
func (s *Session) GetPlayerMap() (*protos.GetMapObjectsResponse, error) {
	cellIDS := s.location.GetCellIDs()
	mapObjRequest, err := proto.Marshal(&protos.GetMapObjectsMessage{
		CellId:           cellIDS,
		SinceTimestampMs: make([]int64, len(cellIDS)),
		Latitude:         s.location.Lat,
		Longitude:        s.location.Lon,
	})
	if err != nil {
		return nil, err
	}
	requests := []*protos.Request{
		{RequestType: protos.RequestType_GET_MAP_OBJECTS, RequestMessage: mapObjRequest},
	}

	response, err := s.Call(requests)
	if err != nil {
		return nil, err
	}

	mapCells := &protos.GetMapObjectsResponse{}
	mapCellBytes := response.Returns[0]
	err = proto.Unmarshal(mapCellBytes, mapCells)
	if err != nil {
		return nil, &ResponseError{err}
	}
	s.feed.Push(mapCells)
	return mapCells, GetErrorFromStatus(response.StatusCode)
}

// GetInventory returns the player items
func (s *Session) GetInventory() (*protos.GetInventoryResponse, error) {
	requests := []*protos.Request{{RequestType: protos.RequestType_GET_INVENTORY}}
	response, err := s.Call(requests)
	if err != nil {
		return nil, err
	}
	inventory := &protos.GetInventoryResponse{}
	err = proto.Unmarshal(response.Returns[0], inventory)
	if err != nil {
		return nil, &ResponseError{err}
	}
	s.feed.Push(inventory)
	return inventory, GetErrorFromStatus(response.StatusCode)
}
