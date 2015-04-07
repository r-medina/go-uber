package uber

import (
	"bytes"
	"fmt"
)

const (
	Version         = "v1"
	RequestEndpoint = "requests"
	ProductEndpoint = "products"
	PriceEndpoint   = "estimates/price"
	TimeEndpoint    = "estimates/time"
	HistoryEndpoint = "history"
	UserEndpoint    = "me"

	// request statuses

	// The `Request` is matching to the most efficient available driver.
	StatusProcessing = "processing"
	// The `Request` was unfulfilled because no drivers were available.
	StatusNoDrivers = "no_drivers_available"
	// The `Request` has been accepted by a driver and is "en route" to the
	// start_location.
	StatusAccepted = "accepted"
	// The driver has arrived or will be shortly.
	StatusArriving = "arriving"
	// The `Request` is "en route" from the start location to the end location.
	StatusInProgress = "in_progress"
	// The `Request` has been canceled by the driver.
	StatusDriverCanceled = "driver_canceled"
	// The `Request` has been canceled by the rider.
	StatusRiderCanceled = "rider_canceled"
	// The `Request` has been completed by the driver.
	StatusCompleted = "completed"

	// the next two use `AUTH_EDPOINT`

	AccessCodeEndpoint  = "authorize"
	AccessTokenEndpoint = "token"

	State = "go-uber"
	Port  = ":7635"
)

// declared as vars so that unit tests can edit the values and hit internal test server
var (
	UberAPIHost = fmt.Sprintf("https://api.uber.com/%s", Version)
	AuthHost    = "https://login.uber.com/oauth"

	UberSandboxAPIHost = fmt.Sprintf("https://sandbox-api.uber.com/%s/sandbox", Version)
)

//
// exported types
//

// Request contains the information relating to a request for an Uber done on behalf of a
// user.
type Request struct {
	RequestID       string `json:"request_id"`
	Status          string `json:"status"`
	Vehicle         `json:"vehicle"`
	Driver          `json:"driver"`
	Location        `json:"location"`
	ETA             int     `json:"eta"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
}

// Vehicle represents the car in a response to requesting a ride.
type Vehicle struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	LicensePlate string `json:"license_plate"`
	PictureURL   string `json:"picture_url"`
}

// Driver represents an Uber driver.
type Driver struct {
	PhoneNumber string `json:"phone_number"`
	Rating      int    `json:"rating"`
	PictureURL  string `json:"picture_url"`
	Name        string `json:"name"`
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

// Price contains information about a price estimate.
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
// given location in seconds.
type Time struct {
	// eg: "5f41547d-805d-4207-a297-51c571cf2a8c"
	ProductID string `json:"product_id"`

	// eg: "UberBLACK"
	DisplayName string `json:"display_name"`

	// The ETA in seconds
	// eg: 410, ie: 6 minutes and 50 seconds
	Estimate int `json:"estimate"`
}

// Location contains a human-readable address as well as the exact coordinates of a location.
type Location struct {
	// Human-readable address
	// eg: "706 Mission St, San Francisco, CA"
	Address string `json:"address,omitempty"`

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
// product type that was requested.
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
// authenticated users profile.
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

//
// internal error types
//

// uberError implements the error interface (by defining an `Error() string` method).
// This datatype is returned from the Uber API with non-2xx responses.
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

// Error implements the `error` interface for `uberError`.
func (err uberError) Error() string {
	var uberErrBuff bytes.Buffer // because O(1) runtime, bitches
	uberErrBuff.WriteString(fmt.Sprintf("Uber API: %s", err.Message))

	// prints code if exists
	if err.Code != "" {
		uberErrBuff.WriteString(fmt.Sprintf("\nCode: %s", err.Code))
	}

	// prints erroneous fields
	if err.Fields != nil {
		uberErrBuff.WriteString("\nFields:")
		for k, v := range err.Fields {
			uberErrBuff.WriteString(fmt.Sprintf("\n\t%s: %v", k, v))
		}
	}

	return uberErrBuff.String()
}
