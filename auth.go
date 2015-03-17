package uber

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

// access is returned from `AccessTokenEndpoint` See comment for `Client.access`.
// https://developer.uber.com/v1/auth/
type access struct {
	// access_token result of three-legged OAuth 2.0 flow. This is needed to make
	// requests on behalf of a user.
	Token string `json:"access_token"`

	// always "Bearer"
	TokenType string `json:"token_type"`

	// 30 days from request
	ExpiresIn int `json:"expires_in"`

	// When the user's access_token has expired, you may obtain a fresh access_token
	// by exchanging the refresh_token that is associated with the access_token
	RefreshToken string `json:"refresh_token"`

	// During the registration process, you'll be prompted to select the scopes your
	// application needs. You can also edit these preferences or add new scopes at a
	// later time. Users will be asked to grant permission for your application's
	// specific scopes when they authenticate.
	//
	// `profile` - Access the basic profile information on a user's Uber account
	// including their first name, email address, and profile picture.
	//
	// `history` - Pull trip data including the locations, times, and product type of
	// a user's historical pickups and drop-offs.
	Scope string `json:"scope"`
}

// auth is the data structure needed to complete OAuth flow.
type auth struct {
	clientID string `query:"client_id,required"`

	// `authReq` uses an `auth`, but it doesn't need a `client_secret`, so we only allow
	// that field in `accReq`
	clientSecret string `query:"-"`
	redirectURI  string `query:"redirect_uri,required"`
}

// OAuth begins the authorization process with Uber. There's no way to do this
// strictly programatically because of the multi-step OAuth process. This method
// returns the URL that the user needs to go to in order for Uber to authorize your
// app and give you a authorization code.
func (c *Client) OAuth(
	clientID, clientSecret, redirect string, scope ...string,
) (string, error) {
	c.auth = new(auth)
	c.clientID = clientID
	c.clientSecret = clientSecret
	c.redirectURI = redirect

	return c.generateRequestURL(AuthHost, AccessCodeEndpoint, authReq{
		auth:         *c.auth,
		responseType: "code",
		scope:        strings.Join(scope, " "), // profile,history
		state:        State,
	})
}

// AutOAuth automatically does the authorization flow by opening the user's browser,
// asking them to authorize, then booting up a server to deal with the user's redirect and
// authorizing your client.
func (c *Client) AutOAuth(
	clientID, clientSecret, redirect string, scope ...string,
) error {
	urlString, err := c.OAuth(clientID, clientSecret, redirect, scope...)
	if err != nil {
		return nil
	}

	httpDone := make(chan struct{})
	httpErr := make(chan error)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != State {
			httpErr <- fmt.Errorf("uber: evidence of tampering--incorrect state %s", state)
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			httpErr <- errors.New("uber: an unidentified error occured")
		}

		err = c.SetAccessToken(code)
		if err != nil {
			httpErr <- err
		}

		fmt.Fprintf(w, `<script type="text/javascript\">close()</script>
you may close this webpage`)
		close(httpDone)
	})

	go func() {
		httpErr <- http.ListenAndServe(Port, nil)
	}()

	err = open.Run(urlString)
	if err != nil {
		return err
	}

	select {
	case err := <-httpErr:
		return err
	case <-httpDone:
		return nil
	}
}

// SetAccessToken completes the third step of the authorization process.
// Once the user generates an authorization code
func (c *Client) SetAccessToken(authorizationCode string) error {
	payload, err := c.generateRequestURLHelper(reflect.ValueOf(accReq{
		auth:         *c.auth,
		clientSecret: c.clientSecret,
		grantType:    "authorization_code",
		code:         authorizationCode,
	}))
	if err != nil {
		return err
	}

	res, err := c.httpClient.PostForm(
		fmt.Sprintf("%s/%s", AuthHost, AccessTokenEndpoint), payload,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode == http.StatusOK {
		access := new(access)
		if err := decoder.Decode(access); err != nil {
			return err
		}

		if access.TokenType == "Bearer" { // always true
			c.access = access
			return nil
		}
	}

	authErr := new(authError)
	decoder.Decode(authErr)
	return authErr
}

// authError is used when there is an error during authentication such that the error
// message can indicate that.
type authError struct {
	// https://developer.uber.com/v1/auth/
	Err string `json:"error"`
}

// Error implements the `error` interface for `authError`.
func (err authError) Error() string {
	return fmt.Sprintf("Authentication: %s", err.Err)
}
