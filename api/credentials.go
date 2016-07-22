package api

// Credentials are used to authenticate against the Pokémon Go API
type Credentials interface {
	GetPassword() string
	GetUsername() string
}

// PTCCredentials are for the Pokémon Trainer's Club
type PTCCredentials struct {
	Username string
	Password string
}

// GetPassword will return the password
func (c *PTCCredentials) GetPassword() string {
	return c.Password
}

// GetUsername returns the username
func (c *PTCCredentials) GetUsername() string {
	return c.Username
}

// GoogleCredentials are for Pokemon Go google accounts
type GoogleCredentials struct {
	Username string
	Password string
}

// GetPassword will return the password
func (c *GoogleCredentials) GetPassword() string {
	return c.Password
}

// GetUsername returns the username
func (c *GoogleCredentials) GetUsername() string {
	return c.Username
}
