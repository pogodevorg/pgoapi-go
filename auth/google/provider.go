package google

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

const androidKeyBase64 = "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="
const androidID = "9774d56d682e549c"
const service = "audience:server:client_id:848232511240-7so421jotr2609rmqakceuu1luuq0ptb.apps.googleusercontent.com"
const app = "com.nianticlabs.pokemongo"
const clientSig = "321187995bc7cdc2b5fc91b11a96e2baa8602c62"

const providerString = "google"

// Provider contains data about and manages the session with the Pokémon Trainer's Club
type Provider struct {
	username string
	password string
	ticket   string
	http     *http.Client
}

// NewProvider constructs a Google auth provider instance
func NewProvider(username, password string) *Provider {

	return &Provider{
		http:     http.DefaultClient,
		username: username,
		password: password,
	}
}

// GetProviderString will return a string identifying the provider
func (p *Provider) GetProviderString() string {
	return providerString
}

// GetAccessToken will return an access token if it has been retrieved
func (p *Provider) GetAccessToken() string {
	return p.ticket
}

// Login retrieves an access token from the Pokémon Trainer's Club
func (p *Provider) Login(ctx context.Context) (string, error) {
	sig, err := signature(p.username, p.password)
	if err != nil {
		return "", err
	}

	postBody := url.Values{}

	postBody.Add("device_country", "us")
	postBody.Add("operatorCountry", "us")
	postBody.Add("lang", "en_US")
	postBody.Add("sdk_version", "23")
	postBody.Add("google_play_services_version", "9256438")
	postBody.Add("accountType", "HOSTED_OR_GOOGLE")
	postBody.Add("Email", p.username)
	postBody.Add("service", service)
	postBody.Add("source", "android")
	postBody.Add("androidId", androidID)
	postBody.Add("app", app)
	postBody.Add("client_sig", clientSig)
	postBody.Add("callerPkg", app)
	postBody.Add("callerSig", clientSig)
	postBody.Add("EncryptedPasswd", sig)

	req, err := http.NewRequest("POST", "https://android.clients.google.com/auth", strings.NewReader(string(postBody.Encode())))
	req.Header.Set("User-Agent", "GoogleAuth/1.4 (mako JDQ39)")
	req.Header.Set("Device", androidID)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("App", app)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	gzBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	decompressedBody, err := ioutil.ReadAll(gzBody)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(decompressedBody), "\n") {
		sp := strings.SplitN(line, "=", 2)
		if len(sp) != 2 {
			continue
		}
		if sp[0] == "Auth" {
			p.ticket = sp[1]
			return p.ticket, nil
		}
	}
	return "", fmt.Errorf("No Auth found")
}

func signature(email, password string) (string, error) {
	androidKeyBytes, err := base64.StdEncoding.DecodeString(androidKeyBase64)
	if err != nil {
		return "", err
	}

	i := bytesToLong(androidKeyBytes[:4]).Int64()
	j := bytesToLong(androidKeyBytes[i+4 : i+8]).Int64()

	androidKey := &rsa.PublicKey{
		N: bytesToLong(androidKeyBytes[4 : 4+i]),
		E: int(bytesToLong(androidKeyBytes[i+8 : i+8+j]).Int64()),
	}

	hash := sha1.Sum(androidKeyBytes)
	msg := append([]byte(email), 0)
	msg = append(msg, []byte(password)...)

	encryptedLogin, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, androidKey, msg, nil)
	if err != nil {
		return "", err
	}

	sig := append([]byte{0}, hash[:4]...)
	sig = append(sig, encryptedLogin...)
	return base64.URLEncoding.EncodeToString(sig), nil
}

func bytesToLong(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}
