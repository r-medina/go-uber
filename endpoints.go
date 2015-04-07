package uber

import "fmt"

//
// the `Client` API
//

// PostRequest allows a ride to be requested on behalf of an Uber user given
// their desired product, start, and end locations.
func (c *Client) PostRequest(
	productID string, startLat, startLon, endLat, endLon float64, surgeConfirmationID string,
) (*Request, error) {
	payload := requestReq{
		productID:           productID,
		startLatitude:       startLat,
		startLongitude:      startLon,
		endLatitude:         endLat,
		endLongitude:        endLon,
		surgeConfirmationID: surgeConfirmationID,
	}
	request := new(requestResp)

	if err := c.httpReqDo("POST", RequestEndpoint, payload, true, request); err != nil {
		return nil, err
	}

	return &request.Request, nil
}

// GetRequest gets the real time status of an ongoing trip that was created using the Ride
// Request endpoint.
func (c *Client) GetRequest(requestID string) (*Request, error) {
	request := new(Request)
	err := c.get(fmt.Sprintf("%s/%s", RequestEndpoint, requestID), nil, true, request)
	if err != nil {
		return nil, err
	}

	return request, nil

}

// DeleteRequest cancels an ongoing `Request` on behalf of a rider.
func (c *Client) DeleteRequest(requestID string) error {
	return c.httpReqDo(
		"DELETE", fmt.Sprintf("%s/%s", RequestEndpoint, requestID), nil, true, nil,
	)
}

// GetRequestMap get a map with a visual representation of a `Request`.
func (c *Client) GetRequestMap(requestID string) (string, error) {
	mapResp := new(requestMapResp)
	err := c.get(fmt.Sprintf("%s/%s/map", RequestEndpoint, requestID), nil, true, mapResp)
	if err != nil {
		return "", err
	}

	return mapResp.HRef, nil
}

// GetProducts returns information about the Uber products offered at a
// given location. The response includes the display name and other details about
// each product, and lists the products in the proper display order.
// https://developer.uber.com/v1/endpoints/#product-types
func (c *Client) GetProducts(lat, lon float64) ([]*Product, error) {
	payload := productsReq{
		latitude:  lat,
		longitude: lon,
	}
	products := new(productsResp)

	if err := c.get(ProductEndpoint, payload, false, products); err != nil {
		return nil, err
	}

	return products.Products, nil
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
		startLongitude: startLon,
		endLatitude:    endLat,
		endLongitude:   endLon,
	}
	prices := new(pricesResp)

	if err := c.get(PriceEndpoint, payload, false, prices); err != nil {
		return nil, err
	}

	return prices.Prices, nil
}

// GetTimes returns ETAs for all products offered at a given location, with the responses
// expressed as integers in seconds. We recommend that this endpoint be called every
// minute to provide the most accurate, up-to-date ETAs.
// The uuid and productID parameters can be empty strings. These provide
// additional experience customization.
func (c *Client) GetTimes(
	startLat, startLon float64, uuid, productID string,
) ([]*Time, error) {
	payload := timesReq{
		startLatitude:  startLat,
		startLongitude: startLon,
		customerUuid:   uuid,
		productID:      productID,
	}
	times := new(timesResp)

	if err := c.get(TimeEndpoint, payload, false, times); err != nil {
		return nil, err
	}

	return times.Times, nil
}

// GetUserActivity returns data about a user's lifetime activity with Uber. The response
// will include pickup locations and times, dropoff locations and times, the distance
// of past requests, and information about which products were requested.
func (c *Client) GetUserActivity(offset, limit int) (*UserActivity, error) {
	payload := historyReq{
		offset: offset,
		limit:  limit,
	}
	userActivity := new(UserActivity)

	if err := c.get(TimeEndpoint, payload, true, userActivity); err != nil {
		return nil, err
	}

	return userActivity, nil
}

// GetUserProfile returns information about the Uber user that has authorized with
// the application.
func (c *Client) GetUserProfile() (*User, error) {
	user := new(User)

	if err := c.get(UserEndpoint, nil, true, user); err != nil {
		return nil, err
	}

	return user, nil
}
