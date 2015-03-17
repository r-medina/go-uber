package uber_test

import (
	"fmt"

	uber "github.com/r-medina/go-uber"
)

var (
	client *uber.Client
)

// In order to use the Uber API you must register an application at
// https://developer.uber.com/. In turn you will receive a `client_id`, `secret`, and
// `server_token`.
func ExampleClient() {
	client = uber.NewClient("your server_token")
}

// All requests require a valid server_token. Requests that require latitude or longitude
// as arguments are float64's (and should be valid coordinates).
func ExampleClient_requests() {
	products, err := client.GetProducts(37.7759792, -122.41823)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, product := range products {
			fmt.Println(*product)
		}
	}

	prices, err := client.GetPrices(41.827896, -71.393034, 41.826025, -71.406892)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, price := range prices {
			fmt.Println(*price)
		}
	}
}

// Uber's OAuth 2.0 flow requires the user go to URL they provide. You can generate this
// URL programatically by calling `client.OAuth`. After the user goes to URL and grants
// your application permissions, you need to figure out a way for the user to input the
// second argument of the URL to which they are redirected (ie: the `code`
// argument`AuthCode` in `RedirectUrl/?state=go-uber&code=AuthCode`). You then need to set
// the access token with `client.SetAccessToken`.
func ExampleClient_auth() {
	url, err := client.OAuth(
		"your client_id", "your client_secret", "your redirect_url", "profile",
		// "profile" is justthe scope. You may also pass another argument "history" (or no
		// scope arguments).
	)
	if err != nil {
		fmt.Printf("OAuth error: %+v\n", err)
		return
	}

	fmt.Printf("Please go to %+v to authorize this app.\n", url)

	client.SetAccessToken("AuthCode")
}

// Alternatively, you can automate the whole thing with the `client.AutoAuth`. This opens
// up a browser for the user.
func ExampleClient_autoAuth() {
	err := client.AutOAuth(
		"your client_id", "your client_secret", "your redirect_url", "profile",
	)

	if err != nil {
		fmt.Printf("OAuth error: %+v\n", err)
	}
}

// After authorizing, you can call the scoped endpoints to which your app has access.
func ExampleClient_userProfile() {
	profile, err := client.GetUserProfile()
	if err != nil {
		fmt.Printf("GetUserProfile error: %+v\n", err)
		return
	}

	fmt.Println(profile)
}
