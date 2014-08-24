/*
The uber package provides an api client for the Uber api.
It exposes methods to get information about Uber products,
estimates, times, and users.

A lot of documentation will be pulled directly from
https://developer.uber.com/v1/endpoints
*/
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

const (
	VERSION          = "v1"
	PRODUCT_ENDPOINT = "products"
	PRICE_ENDPOINT   = "estimates/price"
	TIME_ENDPOINT    = "estimates/time"
	HISTORY_ENDPOINT = "history"
	USER_ENDPOINT    = "me"
)

var UBER_API_ENDPOINT = fmt.Sprintf("http://api.uber.com/%s", VERSION)

// Client stores the tokens needed to access the Uber api.
// All methods of this package that hit said api are methods on this type.
// TODO(asubiott): Abstract the OAuth 2.0 authentication process.
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
	accessToken string
}

// Creates a new client. When accessing a user's profile or activity an
// access token must be specified with the correct scope. If these endpoints
// are not needed, an empty string should be passed in.
func NewClient(serverToken, accessToken string) *Client {
	return &Client{
		serverToken: serverToken,
		accessToken: accessToken,
	}
}

// GetPoducts returns information about the Uber products offered at a
// given location. The response includes the display name and other details about
// each product, and lists the products in the proper display order.
// https://developer.uber.com/v1/endpoints/#product-types
func (c *Client) GetProducts(lat, lon float64) ([]*Product, error) {
	payload := productsReq{
		latitude:  lat,
		longitude: lon,
	}

	products := new([]*Product)
	if err := c.get(PRODUCT_ENDPOINT, payload, products); err != nil {
		return nil, err
	}

	return *products, nil
}

// GetPrices returns an estimated price range for each product offered at a given
// location. The price estimate is provided as a formatted string with the full price
// range and the localized currency symbol.
//
// The response also includes low and high estimates, and the ISO 4217 currency code
// for situations requiring currency conversion. When surge is active for a
// particular product, its surge_multiplier will be greater than 1, but the price
// estimate already factors in this multiplier.
// https://developer.uber.com/v1/endpoints/#price-estimates
func (c *Client) GetPrices(startLat, startLon, endLat, endLon float64) ([]*Price, error) {
	payload := pricesReq{
		startLatitude:  startLat,
		endLatitude:    startLon,
		startLongitude: endLat,
		endLongitude:   endLon,
	}

	prices := new([]*Price)
	if err := c.get(PRICE_ENDPOINT, payload, prices); err != nil {
		return nil, err
	}

	return *prices, nil
}

// GetTimes returns ETAs for all products offered at a given location, with the responses
// expressed as integers in seconds. We recommend that this endpoint be called every
// minute to provide the most accurate, up-to-date ETAs.
// TODO(asubiott): Specify that uuid and productId are optional.
func (c *Client) GetTimes(startLat, startLon float64, uuid, productId string) ([]*Time, error) {
	payload := timesReq{
		startLatitude:  startLat,
		startLongitude: startLon,
		customerUuid:   uuid,
		productId:      productId,
	}

	times := new([]*Time)
	if err := c.get(TIME_ENDPOINT, payload, times); err != nil {
		return nil, err
	}

	return *times, nil
}

// GetHistory returns data about a user's lifetime activity with Uber. The response
// will include pickup locations and times, dropoff locations and times, the distance
// of past requests, and information about which products were requested.
func (c *Client) GetUserActivity(offset, limit int) (*UserActivity, error) {
	payload := historyReq{
		offset: offset,
		limit:  limit,
	}

	userActivity := new(UserActivity)
	if err := c.get(TIME_ENDPOINT, payload, userActivity); err != nil {
		return nil, err
	}

	return userActivity, nil
}

// GetUserProfile returns information about the Uber user that has authorized with
// the application.
func (c *Client) GetUserProfile() (*User, error) {
	payload := userReq{}
	user := new(User)
	if err := c.get(USER_ENDPOINT, payload, user); err != nil {
		return nil, err
	}

	return user, nil
}

// get helps facilitate all the get requests to the Uber api.
// Takes the endpoing, the query parameters, and the data structure that the JSON
// response should be unmarshalled into
func (c *Client) get(endpoint string, payload uberApiRequest, out interface{}) error {
	addr, err := c.generateRequestUrl(endpoint, payload)
	if err != nil {
		return err
	}

	res, err := c.sendRequestWithAuthorization(addr)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(out)
	if err != nil {
		return err
	}

	return nil
}

// sendRequestWithAuthorization sends an HTTP GET request with an Authorization
// field in the header containing the Client's access token (bearer token).
func (c *Client) sendRequestWithAuthorization(url string) (*http.Response, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", c.accessToken))
	return httpClient.Do(req)
}

// generateRequestUrl returns the appropriate a request url to the Uber api based on
// the specified endpoint and the data passed in
func (c *Client) generateRequestUrl(endpoint string, data uberApiRequest) (string, error) {
	payload, err := c.generateRequestUrlHelper(reflect.ValueOf(data))
	if err != nil {
		return "", err
	}

	payload.Add("server_token", c.serverToken)

	// TODO(asubiott): Check if it is harmful in any way to specify the access_token
	// as an empty string and keep the server_token around when we need the access_token.

	return fmt.Sprintf("%s/%s?%s", UBER_API_ENDPOINT, endpoint, payload.Encode()), nil
}

// generateRequestUrlHelper recursively checks `val` to generate the payload. Should
// be used with caution. Only `Client.generateRequestUrl` calls this.
func (c *Client) generateRequestUrlHelper(val reflect.Value) (url.Values, error) {
	payload := make(url.Values)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		queryTag := strings.Split(val.Type().Field(i).Tag.Get("query"), ",")
		// TODO(rm): write a better parser for a query tag--shouldn't just look at
		// position
		if len(queryTag) > 1 && queryTag[1] == "required" {
			if val.Field(i).String() == "" {
				if fieldName == "serverToken" && c.serverToken != "" {
					continue
				}
				return nil, errors.New(fmt.Sprintf("%s is a required field", fieldName))
			}
		}

		var v interface{}
		switch val.Field(i).Kind() {
		case reflect.Int:
			v = val.Field(i).Int()
		case reflect.Float64:
			v = val.Field(i).Float()
		case reflect.String:
			v = val.Field(i).String()
		case reflect.Struct: // we need recursion to support potential future cases
			// where we need to embed more request structs
			supPayload, err := c.generateRequestUrlHelper(val.Field(i))
			if err != nil {
				return nil, err
			}
			// avoids nil field on struct (ie res)
			if len(supPayload) == 0 {
				continue
			}
			for k, va := range supPayload {
				payload.Add(k, va[0])
			}
		default:
			return nil, errors.New(fmt.Sprintf("%s is invalid", fieldName))
		}

		payload.Add(queryTag[0], fmt.Sprintf("%v", v))
	}

	return payload, nil
}
