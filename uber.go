// The uber package provides an api client for the Uber api,
// with methods to get information about Uber products,
// estimates, times, and users.
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

// Creates a new client. When accessing a user's profile or activity an
// access token must be specified with the correct scope. If these endpoints
// are not needed, an empty string should be passed in.
func NewClient(server_token, access_token string) *Client {
	return &Client{
		server_token: server_token,
		access_token: access_token,
	}
}

func (c *Client) GetProducts(lat, lon float64) ([]*Product, error) {
	payload := productsReq{
		latitude:  lat,
		longitude: lon,
	}

	addr, err := c.generateRequestUrl(PRODUCT_ENDPOINT, payload)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	products := new([]*Product)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(products)
	if err != nil {
		return nil, err
	}

	return *products, nil
}

func (c *Client) GetPrices(startLat, startLon, endLat, endLon float64) ([]*Price, error) {
	payload := pricesReq{
		startLatitude:  startLat,
		endLatitude:    startLon,
		startLongitude: endLat,
		endLongitude:   endLon,
	}

	addr, err := c.generateRequestUrl(PRICE_ENDPOINT, payload)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	prices := new([]*Price)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(prices)
	if err != nil {
		return nil, err
	}

	return *prices, nil
}

func (c *Client) GetTimes(startLat, startLon float64, uuid, productId string) ([]*Time, error) {
	payload := timesReq{
		startLatitude:  startLat,
		startLongitude: startLon,
		customerUuid:   uuid,
		productId:      productId,
	}

	addr, err := c.generateRequestUrl(TIME_ENDPOINT, payload)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	times := new([]*Time)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(times)
	if err != nil {
		return nil, err
	}

	return *times, nil
}

func (c *Client) GetUserProfile() (*User, error) {
	addr, err := c.generateRequestUrl(USER_ENDPOINT, nil)
	if err != nil {
		return nil, err
	}

	// TODO(asubiott): Should there be a check for correct access_token scope?
	// This might be taken care of for us by the Uber API. It should return
	// an error.

	res, err := c.sendRequestWithAuthorization(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	user := new(User)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Sends an HTTP GET request with an Authorization field in the header
// containing the Client's access token.
func (c *Client) sendRequestWithAuthorization(url string) (*http.Response, error) {
	http_client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", c.access_token))
	return http_client.Do(req)
}

func (c *Client) generateRequestUrl(endpoint string, data interface{}) (string, error) {
	payload, err := c.generateRequestUrlHelper(reflect.ValueOf(data))
	if err != nil {
		return "", err
	}

	payload.Add("server_token", c.server_token)
	// TODO(asubiott): Check if it is harmful in any way to specify the access_token
	// as an empty string and keep the server_token around when we need the access_token.

	return fmt.Sprintf("%s/%s?%s", UBER_API_ENDPOINT, endpoint, payload.Encode()), nil
}

func (c *Client) generateRequestUrlHelper(val reflect.Value) (url.Values, error) {
	payload := make(url.Values)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		queryTag := strings.Split(val.Type().Field(i).Tag.Get("query"), ",")
		if len(queryTag) > 1 && queryTag[1] == "required" {
			if val.Field(i).String() == "" {
				if fieldName == "serverToken" && c.server_token != "" {
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
