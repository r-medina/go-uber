package uber

// access is returned from `AccessTokenEndpoint`
// See comment for `Client.access`
// https://developer.uber.com/v1/auth/
type access struct {
	// access_token result of three-legged OAuth 2.0 flow
	// Needed to make requests on behalf of a user
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

// Product type specifies an Uber product.
// An Uber product refers to a specific type of car/service.
type Product struct {
	// Unique identifier representing a specific product for a given latitude &
	// longitude. For example, uberX in San Francisco will have a different
	// product_id than uberX in Los Angeles.
	// eg: "327f7914-cd12-4f77-9e0c-b27bac580d03"
	ProductID string `json:"product_id"`

	// Description of product
	// eg: "The original Uber"
	Description string `json:"description"`

	// eg: "UberBLACK"
	DisplayName string `json:"display_name"`

	// eg: 4
	Capacity int `json:"capacity"`

	// A URI specifying the location of an image
	// eg: "http://..."
	Image string `json:"image"`
}

// productsResp is the type that is returned from the `ProductEndpoint`
// This data definition is needed so that unmarshalling can actually happen
type productsResp struct {
	Products []*Product `json:"products"`
}

// Price contains information about a price estimate
type Price struct {
	// eg: "08f17084-23fd-4103-aa3e-9b660223934b"
	ProductID string `json:"product_id"`

	// ISO 4217 currency code for situations requiring currency conversion
	// eg: "USD"
	CurrencyCode string `json:"currency_code"`

	// eg: "UberBLACK"
	DisplayName string `json:"display_name"`

	// Formatted string of estimate in local currency of the start location. Estimate
	// could be a range, a single number (flat rate) or "Metered" for TAXI.
	// eg: "$23-29"
	Estimate string `json:"estimate"`

	// The lowest value in the estimate for the given currency
	// eg: 23
	LowEstimate int `json:"low_estimate,string"`

	// The highest value in the estimate for the given currency
	// eg: 29
	HighEstimate int `json:"high_estimate,string"`

	// Uber price gouging factor
	// http://www.technologyreview.com/review/529961/in-praise-of-efficient-price-gouging/
	// eg: 1
	SurgeMultiplier float64 `json:"surge_multiplier"`
}

// pricesResp is the type that is returned from the `PriceEndpoint`
// This data definition is needed so that unmarshalling can actually happen
type pricesResp struct {
	Prices []*Price `json:"prices"`
}

// Time contains information about the estimated time of arrival for a product at a
// given location in seconds
type Time struct {
	// eg: "5f41547d-805d-4207-a297-51c571cf2a8c"
	ProductID string `json:"product_id"`

	// eg: "UberBLACK"
	DisplayName string `json:"display_name"`

	// The ETA in seconds
	// eg: 410, ie: 6 minutes and 50 seconds
	Estimate int `json:"estimate"`
}

// pricesResp is the type that is returned from the `PriceEndpoint`
// This data definition is needed so that unmarshalling can actually happen
type timesResp struct {
	Times []*Time `json:"times"`
}

// Location contains a human-readable address as well as the exact coordinates of a location
type Location struct {
	// Human-readable address
	// eg: "706 Mission St, San Francisco, CA"
	Address string `json:"address"`

	// eg: 37.7860099
	Latitude float64 `json:"latitude"`

	// eg: -122.4025387
	Longitude float64 `json:"longitude"`
}

// UserActivity contains data about a user's lifetime activity with Uber.
type UserActivity struct {
	// How much the list of returned results is offset by (position in pagination)
	// eg: 0
	Offset int `json:"offset"`

	// Number of items retrieved (that is, the length of `History` in this struct,
	// but not the total length of the history)
	// eg: 1
	Limit int `json:"limit"`

	// Total number of items available
	// eg: 5
	Count int `json:"count"`

	// List of trips (see `Trip`)
	History []*Trip `json:"history"`
}

// Trip contains Information including the pickup location, dropoff location, request
// start time, request end time, and distance of requests (in miles), as well as the
// product type that was requested
type Trip struct {
	// Customer UUID
	// eg: "7354db54-cc9b-4961-81f2-0094b8e2d215"
	Uuid string `json:"uuid"`

	// Time in seconds
	// eg: 1401884467
	RequestTime int `json:"request_time"`

	// eg: edf5e5eb-6ae6-44af-bec6-5bdcf1e3ed2c
	ProductID string `json:"product_id"`

	// String depicting the status of the trip. Don't know what values these could take
	// because the website only shows "completed"
	// eg: "completed"
	Status string `json:"status"`

	// Distance of request in miles (presumable that of the customer to he nearest driver)
	// eg: 0.0279562
	Distance float64 `json:"distance"`

	// Start time of trip
	// eg: 1401884646
	StartTime int `json:"start_time"`

	// Self explanatory (see `Location`)
	StartLocation *Location `json:"start_location"`

	// Start time of trip
	// eg: 1401884732
	EndTime int `json:"end_time"`

	// Self explanatory (see `Location`)
	EndLocation *Location `json:"end_location"`
}

// User is the response from the /me endpoint. Provides information about the
// authenticated users profile
type User struct {
	// eg: "Uber"
	FirstName string `json:"first_name"`

	// eg: "Developer"
	LastName string `json:"last_name"`

	// eg: "developer@uber.com"
	Email string `json:"email"`

	// Image URI
	// eg: "https://..."
	Picture string `json:"picture"`

	// Promotion code user has activated
	// eg: "teypo"
	PromoCode string `json:"promo_code"`
}

// TODO(r-medina): add doc
type auth struct {
	clientID     string `query:"client_id,required"`
	clientSecret string `query:"-"`
	redirectURI  string `query:"redirect_uri,required"`
}

// TODO(r-medina): add doc
type authReq struct {
	auth
	responseType string `query:"response_type,required"`
	scope        string `query:"scope"`
	state        string `query:"state"`
}

// TODO(r-medina): add doc
type accReq struct {
	auth
	clientSecret string `query:"client_secret,required"`
	grantType    string `query:"grant_type,required"`
	code         string `query:"code,required"`
}

// TODO(r-medina): add doc
type productsReq struct {
	latitude  float64 `query:"latitude,required"`
	longitude float64 `query:"longitude,required"`
}

// TODO(r-medina): add doc
type pricesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	endLatitude    float64 `query:"end_latitude,required"`
	endLongitude   float64 `query:"end_longitude,required"`
}

// TODO(r-medina): add doc
type timesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	customerUuid   string  `query:"customer_uuid"`
	productID      string  `query:"product_id"`
}

// TODO(r-medina): add doc
type historyReq struct {
	offset int `query:"offset,required"`
	limit  int `query:"limit,required"`
}

// uberError implements the error interface (by defining an `Error() string` method).
// This datatype is returned from the Uber api with non-2xx responses
type uberError struct {
	// Human readable message which corresponds to the client error
	// eg: "Invalid user"
	Message string `json:"message"`

	// Underscored delimited string
	// eg: "invalid"
	Code string `json:"code"`

	// A hash of field names that have validations. This has a value of an array with
	// member strings that describe the specific validation error
	// eg: map{"first_name": ["Required"]}
	Fields map[string]string `json:"fields,omitempty"`
}

// TODO(r-medina): add doc
type authError struct {
	// https://developer.uber.com/v1/auth/
	error string `json:"error"`
}
