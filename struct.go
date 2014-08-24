package uber

// Product type specifies an Uber product.
// An Uber product refers to a specific type of car/service.
type Product struct {
	// Unique identifier representing a specific product for a given latitude &
	// longitude. For example, uberX in San Francisco will have a different
	// product_id than uberX in Los Angeles.
	// eg: "327f7914-cd12-4f77-9e0c-b27bac580d03"
	ProductId string `json:"product_id"`

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

// Price contains information about a price estimate
type Price struct {
	// eg: "08f17084-23fd-4103-aa3e-9b660223934b"
	ProductId string `json:"product_id"`

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
	LowEstimate int `json:"low_estimate"`

	// The highest value in the estimate for the given currency
	// eg: 29
	HighEstimate int `json:"high_estimate"`

	// Uber price gouging factor
	// http://www.technologyreview.com/review/529961/in-praise-of-efficient-price-gouging/
	// eg: 1
	SurgeMultiplier float64 `json:"surge_multiplier"`
}

// Time contains information about the estimated time of arrival for a product at a
// given location in seconds
type Time struct {
	// eg: "5f41547d-805d-4207-a297-51c571cf2a8c"
	ProductId string `json:"product_id"`

	// eg: "UberBLACK"
	DisplayName string `json:"display_name"`

	// The ETA in seconds
	// eg: 410, ie: 6 minutes and 50 seconds
	Estimate int `json:"estimate"`
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
	// TODO(rm): find out more about this
	// eg: 1401884467
	RequestTime int `json:"request_time"`

	// eg: edf5e5eb-6ae6-44af-bec6-5bdcf1e3ed2c
	ProductId string `json:"product_id"`

	// String depicting the status of the trip. Don't know what values these could take
	// because the website only shows "completed"
	// eg: "completed"
	Status string `json:"status"`

	// Distance of request in miles (presumable that of the customer to he nearest driver)
	// eg: 0.0279562
	Distance float64 `json:"distance"`

	// Start time of trip
	// See TODO for `RequestTime`
	// eg: 1401884646
	StartTime int `json:"start_time"`

	// Self explanatory (see `Location`)
	StartLocation *Location `json:"start_location"`

	// Start time of trip
	// See TODO for `RequestTime`
	// eg: 1401884732
	EndTime int `json:"end_time"`

	// Self explanatory (see `Location`)
	EndLocation *Location `json:"end_location"`
}

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

// uberApiRequest is a shell data definition that is just used to document that
// `Client.generateRequestUrl` takes a specific type of data
type uberApiRequest interface{}

type productsReq struct {
	latitude  float64 `query:"latitude,required"`
	longitude float64 `query:"longitude,required"`
}

type pricesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	endLatitude    float64 `query:"end_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	endLongitude   float64 `query:"end_longitude,required"`
}

type timesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	customerUuid   string  `query:"customer_uuid"`
	productId      string  `query:"product_id"`
}

type historyReq struct {
	offset int `query:"offset,required"`
	limit  int `query:"limit,required"`
}

type userReq struct{}
