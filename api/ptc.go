package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const ptcAuthorizeURL = "https://sso.pokemon.com/sso/oauth2.0/accessToken"
const ptcLoginURL = "https://sso.pokemon.com/sso/login?service=https://sso.pokemon.com/sso/oauth2.0/callbackAuthorize"

const pctRedirectURI = "https://www.nianticlabs.com/pokemongo/error"
const pctClientSecret = "w8ScCUXJQc6kXKw8FiOhd8Fixzht18Dq3PEVkUCP5ZPxtgyWsbTvWHFLm2wNY0JR"
const pctClientID = "mobile-app_pokemon-go"

type ptcLoginRequest struct {
	Lt        string   `json:"lt"`
	Execution string   `json:"execution"`
	Errors    []string `json:"errors,omitempty"`
}

// LoginPTCError is thrown when something went wrong with the login request
type LoginPTCError struct {
	message string
}

func (e *LoginPTCError) Error() string {
	return fmt.Sprintf("login/ptc: %s", e.message)
}

// getHTTPClient returns a client that has a cookie jar to save sessions for subsequent requests
func getHTTPClient() *http.Client {
	options := &cookiejar.Options{}
	jar, _ := cookiejar.New(options)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Use the last error")
		},
	}
	return client
}

func getLoginPTCError(message string) (string, error) {
	return "", &LoginPTCError{message}
}

// getPTCAccessToken retrieves the access token for Pokémon Trainer's Club
func getPTCAccessToken(user, pass string) (accessToken string, err error) {
	client := getHTTPClient()

	req1, _ := http.NewRequest("GET", ptcLoginURL, nil)
	req1.Header.Set("User-Agent", "niantic")

	resp1, err1 := client.Do(req1)
	if err1 != nil {
		return getLoginPTCError("Could not start login process, the website might be down")
	}

	defer resp1.Body.Close()
	body1, _ := ioutil.ReadAll(resp1.Body)
	var loginRespBody ptcLoginRequest
	json.Unmarshal(body1, &loginRespBody)
	resp1.Body.Close()

	loginForm := url.Values{}
	loginForm.Set("lt", loginRespBody.Lt)
	loginForm.Set("execution", loginRespBody.Execution)
	loginForm.Set("_eventId", "submit")
	loginForm.Set("username", user)
	loginForm.Set("password", pass)

	loginFormData := strings.NewReader(loginForm.Encode())

	req2, _ := http.NewRequest("POST", ptcLoginURL, loginFormData)
	req2.Header.Set("User-Agent", "niantic")
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, err2 := client.Do(req2)
	if _, ok2 := err2.(*url.Error); !ok2 {

		defer resp2.Body.Close()
		body2, _ := ioutil.ReadAll(resp2.Body)
		var respBody ptcLoginRequest
		json.Unmarshal(body2, &respBody)
		resp2.Body.Close()

		if len(respBody.Errors) > 0 {
			return getLoginPTCError(respBody.Errors[0])
		}

		return getLoginPTCError("Could not request authorization")
	}

	location, _ := url.Parse(resp2.Header.Get("Location"))
	ticket := location.Query().Get("ticket")

	authorizeForm := url.Values{}
	authorizeForm.Set("client_id", pctClientID)
	authorizeForm.Set("redirect_uri", pctRedirectURI)
	authorizeForm.Set("client_secret", pctClientSecret)
	authorizeForm.Set("grant_type", "refresh_token")
	authorizeForm.Set("code", ticket)

	authorizeFormData := strings.NewReader(authorizeForm.Encode())

	req3, _ := http.NewRequest("POST", ptcAuthorizeURL, authorizeFormData)
	req2.Header.Set("User-Agent", "niantic")
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp3, err3 := client.Do(req3)
	if err3 != nil {
		return getLoginPTCError("Could not authorize code")
	}

	b, _ := ioutil.ReadAll(resp3.Body)
	query, _ := url.ParseQuery(string(b))

	return query.Get("access_token"), nil
}
