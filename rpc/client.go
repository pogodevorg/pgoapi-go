package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/golang/protobuf/proto"

	"github.com/pkmngo-odi/pogo-protos"
)

const httpUserAgent = "Niantic App"

func raise(message string) error {
	return fmt.Errorf("rpc/client: %s", message)
}

// Client is used to communicate with the Pokémon Go API
type Client struct {
	http *http.Client
}

// NewClient constructs a Pokémon Go RPC API client
func NewClient() *Client {
	options := &cookiejar.Options{}
	jar, _ := cookiejar.New(options)
	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return raise("Did not follow redirect")
		},
	}

	return &Client{
		http: httpClient,
	}
}

// Request queries the Pokémon Go API will all pending requests
func (c *Client) Request(endpoint string, requestEnvelope *protos.RequestEnvelope) (responseEnvelope *protos.ResponseEnvelope, err error) {
	responseEnvelope = &protos.ResponseEnvelope{}

	// Build request
	requestBytes, err := proto.Marshal(requestEnvelope)
	if err != nil {
		return responseEnvelope, raise("Could not encode request body")
	}
	requestReader := bytes.NewReader(requestBytes)
	request, err := http.NewRequest("POST", endpoint, requestReader)
	if err != nil {
		return responseEnvelope, raise("Unable to create the request")
	}
	request.Header.Add("User-Agent", httpUserAgent)

	// Perform call to API
	response, err := c.http.Do(request)
	if err != nil {
		return responseEnvelope, raise(fmt.Sprintf("There was an error requesting the API: %s", err))
	}
	if response.StatusCode != 200 {
		return responseEnvelope, raise(fmt.Sprintf("Status code was %d, expected 200", response.StatusCode))
	}

	// Read the response
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseEnvelope, raise("Could not decode response body")
	}

	proto.Unmarshal(responseBytes, responseEnvelope)

	return responseEnvelope, nil
}
