package api

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	protos "github.com/pogodevorg/POGOProtos-go"
	"github.com/pogodevorg/pgoapi-go/auth"
	"gitlab.com/pidgeyfinder/pokecrypt"
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

	hasTicket bool
	ticket    *protos.AuthTicket
	started   time.Time
	provider  auth.Provider
	hash      []byte
}

func generateRequests() []*protos.Request {
	return make([]*protos.Request, 0)
}

func getTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(time.Millisecond))
}

// NewSession constructs a Pokémon Go RPC API client
func NewSession(provider auth.Provider, location *Location, feed Feed, debug bool) *Session {
	return &Session{
		location:  location,
		rpc:       NewRPC(),
		provider:  provider,
		debug:     debug,
		debugger:  &jsonpb.Marshaler{Indent: "\t"},
		feed:      feed,
		started:   time.Now(),
		hasTicket: false,
		hash:      make([]byte, 32),
	}
}

// SetTimeout sets the client timeout for the RPC API
func (s *Session) SetTimeout(d time.Duration) {
	s.rpc.http.Timeout = d
}

func (s *Session) setTicket(ticket *protos.AuthTicket) {
	s.hasTicket = true
	s.ticket = ticket
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

func (s *Session) debugProtoMessage(label string, pb proto.Message) {
	if s.debug {
		str, _ := s.debugger.MarshalToString(pb)
		log.Println(fmt.Sprintf("%s: %s", label, str))
	}
}

// Call queries the Pokémon Go API through RPC protobuf
func (s *Session) Call(ctx context.Context, requests []*protos.Request) (*protos.ResponseEnvelope, error) {

	requestEnvelope := &protos.RequestEnvelope{
		RequestId:  uint64(8145806132888207460),
		StatusCode: int32(2),

		MsSinceLastLocationfix: int64(989),

		Longitude: s.location.Lon,
		Latitude:  s.location.Lat,

		Accuracy: s.location.Accuracy,

		Requests: requests,
	}

	if s.hasTicket {
		requestEnvelope.AuthTicket = s.ticket
	} else {
		requestEnvelope.AuthInfo = &protos.RequestEnvelope_AuthInfo{
			Provider: s.provider.GetProviderString(),
			Token: &protos.RequestEnvelope_AuthInfo_JWT{
				Contents: s.provider.GetAccessToken(),
				Unknown2: int32(59),
			},
		}
	}

	if s.hasTicket {
		t := getTimestamp(time.Now())
		ticket, _ := proto.Marshal(s.ticket)
		requestHash := make([]uint64, len(requests))

		for idx, request := range requests {
			req, err := proto.Marshal(request)
			if err != nil {
				return nil, err
			}
			requestHash[idx] = pokecrypt.HashRequest(ticket, req)
		}

		locationHash1 := pokecrypt.HashLoaction1(ticket, s.location.Lat, s.location.Lon, s.location.Alt)
		locationHash2 := pokecrypt.HashLocation2(s.location.Lat, s.location.Lon, s.location.Alt)

		signature := &protos.Signature{
			RequestHash:   requestHash,
			LocationHash1: int32(locationHash1),
			LocationHash2: int32(locationHash2),
			ActivityStatus: &protos.Signature_ActivityStatus{
				Stationary: true,
			},
			DeviceInfo: &protos.Signature_DeviceInfo{
				DeviceId:             "<device_id>",
				DeviceBrand:          "Apple",
				DeviceModel:          "iPhone",
				DeviceModelBoot:      "Iphone7,2",
				HardwareManufacturer: "Apple",
				HardwareModel:        "N66AP",
				FirmwareBrand:        "iPhone OS",
				FirmwareType:         "9.3.3",
			},
			SessionHash:         s.hash,
			Timestamp:           t,
			TimestampSinceStart: (t - getTimestamp(s.started)),
			Unknown25:           pokecrypt.Hash25(),
		}

		signatureProto, err := proto.Marshal(signature)
		if err != nil {
			return nil, ErrFormatting
		}

		requestMessage, err := proto.Marshal(&protos.SendEncryptedSignatureRequest{
			EncryptedSignature: pokecrypt.Encrypt(signatureProto, uint32(signature.TimestampSinceStart)),
		})
		if err != nil {
			return nil, ErrFormatting
		}

		requestEnvelope.PlatformRequests = []*protos.RequestEnvelope_PlatformRequest{
			{
				Type:           protos.PlatformRequestType_SEND_ENCRYPTED_SIGNATURE,
				RequestMessage: requestMessage,
			},
		}

		s.debugProtoMessage("request signature", signature)
	}

	s.debugProtoMessage("request envelope", requestEnvelope)

	responseEnvelope, err := s.rpc.Request(ctx, s.getURL(), requestEnvelope)

	s.debugProtoMessage("response envelope", responseEnvelope)

	return responseEnvelope, err
}

// MoveTo sets your current location
func (s *Session) MoveTo(location *Location) {
	s.location = location
}

// Init initializes the client by performing full authentication
func (s *Session) Init(ctx context.Context) error {

	var err error;
	tokenExists := s.provider.AccessTokenExists()

  switch tokenExists {
	case false:
		_, err = s.provider.Login(ctx)
		if err != nil {
			return err
		}
	default:
		break;
	}

	_, err = rand.Read(s.hash)
	if err != nil {
		return ErrFormatting
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

	response, err := s.Call(ctx, requests)
	if err != nil {
		return err
	}

	url := response.ApiUrl
	if url == "" {
		return ErrNoURL
	}
	s.setURL(url)

	ticket := response.GetAuthTicket()
	s.setTicket(ticket)

	return nil
}

// Announce publishes the player's presence and returns the map environment
func (s *Session) Announce(ctx context.Context) (mapObjects *protos.GetMapObjectsResponse, err error) {

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
		{RequestType: protos.RequestType_CHECK_CHALLENGE},
	}

	response, err := s.Call(ctx, requests)
	if err != nil {
		return mapObjects, ErrRequest
	}

	mapObjects = &protos.GetMapObjectsResponse{}
	if len(response.Returns) < 5 {
		return nil, errors.New("Empty response")
	}
	err = proto.Unmarshal(response.Returns[5], mapObjects)
	if err != nil {
		return nil, &ErrResponse{err}
	}
	s.feed.Push(mapObjects)
	s.debugProtoMessage("response return[5]", mapObjects)

	return mapObjects, GetErrorFromStatus(response.StatusCode)
}

// GetPlayer returns the current player profile
func (s *Session) GetPlayer(ctx context.Context) (*protos.GetPlayerResponse, error) {
	requests := []*protos.Request{{RequestType: protos.RequestType_GET_PLAYER}}
	response, err := s.Call(ctx, requests)
	if err != nil {
		return nil, err
	}

	player := &protos.GetPlayerResponse{}
	err = proto.Unmarshal(response.Returns[0], player)
	if err != nil {
		return nil, &ErrResponse{err}
	}
	s.feed.Push(player)
	s.debugProtoMessage("response return[0]", player)

	return player, GetErrorFromStatus(response.StatusCode)
}

// GetPlayerMap returns the surrounding map cells
func (s *Session) GetPlayerMap(ctx context.Context) (*protos.GetMapObjectsResponse, error) {
	return s.Announce(ctx)
}

// GetInventory returns the player items
func (s *Session) GetInventory(ctx context.Context) (*protos.GetInventoryResponse, error) {
	requests := []*protos.Request{{RequestType: protos.RequestType_GET_INVENTORY}}
	response, err := s.Call(ctx, requests)
	if err != nil {
		return nil, err
	}
	inventory := &protos.GetInventoryResponse{}
	err = proto.Unmarshal(response.Returns[0], inventory)
	if err != nil {
		return nil, &ErrResponse{err}
	}
	s.feed.Push(inventory)
	s.debugProtoMessage("response return[0]", inventory)

	return inventory, GetErrorFromStatus(response.StatusCode)
}
