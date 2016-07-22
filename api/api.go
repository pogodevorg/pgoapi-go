package api

// Client is used to communicate with the Pokémon Go API
type Client struct {
	credentials Credentials
	accessToken string
}

// NewClient creates a new API client with a set of credentials
func NewClient(credentials Credentials) *Client {
	return &Client{
		credentials: credentials,
	}
}

// GetAccessToken retrieves the access token from the credentials
func (c *Client) GetAccessToken() (accessToken string, err error) {
	if c.accessToken != "" {
		return c.accessToken, nil
	}

	if _, ok := c.credentials.(*PTCCredentials); ok {
		accessToken, err = getPTCAccessToken(c.credentials.GetUsername(), c.credentials.GetPassword())
	}

	c.accessToken = accessToken
	return c.accessToken, err
}
