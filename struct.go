package uber

// Client stores the tokens needed to access the Uber api.
// All methods of this package that hit said api are methods on this type.
type Client struct {
	// TODO(asubiott): document these
	server_token string
	access_token string
}

// Product type specifies an Uber product.
// An Uber product refers to a specific type of car/service.
type Product struct {
	// eg: "327f7914-cd12-4f77-9e0c-b27bac580d03"
	ProductId string `json:"product_id"`

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

	// A human-readable price estimate
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

type UserHistory struct {
	Offset  int     `json:"offset"`
	Limit   int     `json:"limit"`
	Count   int     `json:"count"`
	History []*Trip `json:"history"`
}

type Trip struct {
	Uuid          string    `json:"uuid"`
	RequestTime   int       `json:"request_time"`
	ProductId     string    `json:"product_id"`
	Status        string    `json:"status"`
	Distance      float64   `json:"distance"`
	StartTime     int       `json:"start_time"`
	StartLocation *Location `json:"start_location"`
	EndTime       int       `json:"end_time"`
	EndLocation   *Location `json:"end_location"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	PromoCode string `json:"promo_code"`
}

type req struct {
	serverToken string `query:"server_token,required"`
}

type productsReq struct {
	req
	latitude  float64 `query:"latitude,required"`
	longitude float64 `query:"longitude,required"`
}

type pricesReq struct {
	req
	startLatitude  float64 `query:"start_latitude,required"`
	endLatitude    float64 `query:"end_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	endLongitude   float64 `query:"end_longitude,required"`
}

type timesReq struct {
	req
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	customerUuid   string  `query:"customer_uuid,required"`
	productId      string  `query:"product_id,required"`
}

type historyReq struct {
	offset int `query:"offset, required"`
	limit  int `query:"limit, required"`
}
