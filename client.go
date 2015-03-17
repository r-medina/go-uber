package uber

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Client stores the tokens needed to access the Uber api.
// All methods of this package that hit said api are methods on this type.
type Client struct {
	// Your API token should be specified if your application will access the
	// Products, Price Estimates, and Time Estimates endpoints.
	serverToken string

	// OAuth 2.0 bearer token necessary for the use of the User Activity and
	// User Profile endpoints. It is the result of three step authentication
	// outlined in https://developer.uber.com/v1/auth/#oauth-2-0. When procuring
	// this token, keep in mind that you must specify the history scope if you
	// intend to use the User Activity endpoint and the profile scope if you
	// intend to use the User Profile endpoint.
	*access

	// An http.Client is needed to make requests to the API as well as do the
	// authentication. Rather than instantiate a new client on each request, we
	// memoize it here, as it will always be used.
	httpClient *http.Client

	// contains further authentication information for Uber OAuth flow.
	*auth
}

// NewClient creates a new client. The serverToken is your API token provided by Uber.
// When accessing a user's profile or activity a serverToken is not enough and an
// accessToken must be specified with the correct scope.
// To access those endpoints, use `*Client.OAuth()`
func NewClient(serverToken string) *Client {
	return &Client{
		serverToken: serverToken,
		access:      new(access),
		httpClient:  new(http.Client),
	}
}

// get helps facilitate all the get requests to the Uber api.
// Takes the endpoint, the query parameters, whether or not oauth should be used
// and the data structure that the JSON response should be unmarshalled into.
func (c *Client) get(endpoint string, payload uberAPIReq, oauth bool, out uberAPIResp) error {
	return c.httpReqDo("GET", endpoint, payload, oauth, out)
}

func (c *Client) httpReqDo(
	method, endpoint string, payload uberAPIReq, oauth bool, out uberAPIResp,
) error {
	url, err := c.generateRequestURL(UberAPIHost, endpoint, payload)
	if err != nil {
		return err
	}

	res, err := c.sendRequestWithAuthorization(method, url, oauth)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	// If the status code is non-2xx, generate the error
	switch {
	case res.StatusCode == http.StatusNotFound:
		// should never, ever happen because we specify the endpoints
		return &uberError{
			Message: fmt.Sprintf("Endpoint '%s' not found.", endpoint),
		}
	case res.StatusCode >= 300:
		decoder = json.NewDecoder(res.Body)

		// no good way to do this with `http.Status...` codes ;o
		uberErr := new(uberError)
		if err := decoder.Decode(uberErr); err != nil {
			return err
		}

		// the case where the Uber api didn't provide an UberError in the response
		if uberErr.Message == "" && uberErr.Code == "" {
			return errors.New("uber: an unidentified error occured")
		}

		return *uberErr
	}

	err = decoder.Decode(out)
	if err != nil {
		return err
	}

	return nil
}

// sendRequestWithAuthorization sends an HTTP GET request with an Authorization
// field in the header containing the Client's access token (bearer token) if
// the oauth parameter is true and the server token (api token) if not.
func (c *Client) sendRequestWithAuthorization(
	method, url string, oauth bool,
) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	authStr := fmt.Sprintf("Token %s", c.serverToken)
	if oauth {
		authStr = fmt.Sprintf("Bearer %s", c.Token)
	}

	req.Header.Set("authorization", authStr)

	return c.httpClient.Do(req)
}

// generateRequestURL returns the appropriate a request url to the Uber api based on
// the specified endpoint and the data passed in
func (c *Client) generateRequestURL(base, endpoint string, data uberAPIReq) (string, error) {
	var queryParameters string
	if data == nil {
		queryParameters = ""
	} else {
		payload, err := c.generateRequestURLHelper(reflect.ValueOf(data))
		if err != nil {
			return "", err
		}

		queryParameters = payload.Encode()
	}

	if queryParameters != "" {
		queryParameters = fmt.Sprintf("?%s", queryParameters)
	}

	return fmt.Sprintf("%s/%s%s", base, endpoint, queryParameters), nil
}

// generateRequestURLHelper recursively checks `val` to generate the payload. Should
// be used with caution. Only `Client.generateRequestURL` calls this.
func (c *Client) generateRequestURLHelper(val reflect.Value) (url.Values, error) {
	payload := make(url.Values)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		queryTag := strings.Split(val.Type().Field(i).Tag.Get("query"), ",")

		if queryTag[0] == "-" { // skip this field
			continue
		}

		var v interface{}
		switch val.Field(i).Kind() {
		case reflect.Int:
			v = val.Field(i).Int()
		case reflect.Float64:
			v = val.Field(i).Float()
		case reflect.String:
			v = val.Field(i).String()
			if len(queryTag) > 1 && queryTag[1] == "required" {
				// cannot be required and empty
				if v == "" {
					return nil, fmt.Errorf("uber: %s is a required field", fieldName)
				}
			}
		case reflect.Struct:
			supPayload, err := c.generateRequestURLHelper(val.Field(i))
			if err != nil {
				return nil, err
			}

			// avoids nil field on struct (eg res)
			if len(supPayload) == 0 {
				continue
			}

			for k, va := range supPayload {
				payload.Add(k, va[0])
			}
		default:
			return nil, fmt.Errorf("%s is invalid", fieldName)
		}

		if v != "" && queryTag[0] != "" {
			payload.Add(queryTag[0], fmt.Sprintf("%v", v))
		}
	}

	return payload, nil
}

// Shell data definitions used to document that `Client.generateRequestURL` takes a
// specific type of data
type uberAPIReq interface{}
type uberAPIResp interface{}
