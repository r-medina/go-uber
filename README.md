go-uber
=======

Golang bindings for [Uber API](https://developer.uber.com/v1/endpoints/)

[![Build Status](https://travis-ci.org/r-medina/go-uber.png)](https://travis-ci.org/r-medina/go-uber)
[![GoDoc](https://godoc.org/github.com/r-medina/go-uber?status.svg)](https://godoc.org/github.com/r-medina/go-uber)

# Usage

## Register Your Application

In order to use the Uber API you must register an application at the [Uber Developer Portal](https://developer.uber.com).
In turn you will receive a `client_id`, `secret`, and `server_token`.

## Creating a Client

```go
package main

import (
  "github.com/r-medina/go-uber"
)

func main() {
	client := uber.NewClient(SERVER_TOKEN)
}
```

## Making Requests

Currently, the Uber API offers support for requesting information about products (e.g available cars), price estimates, time estimates, user ride history, and user info. All requests require a valid `server_token`. Requests that require latitude or longitude as arguments are float64's (and should be valid lat/lon's).

For more information about using this library, view the [Godocs](http://godoc.org/github.com/r-medina/go-uber).

```go
products, err := c.GetProducts(37.7759792, -122.41823)
if err != nil {
	fmt.Println(err)
} else {
	for _, product := range productss {
		fmt.Println(*product)
	}
}

prices, err := c.GetPrices(41.827896, 41.826025, -71.393034, -71.406892)
if err != nil {
	fmt.Println(err)
} else {
	for _, price := range prices {
		fmt.Println(*price)
	}
}
```
## Authorizing

Uber's OAuth 2.0 flow requires the user go to URL they provide. As of right now, after the user authorizes the app, there is no elegant way to get it into the client, but we're working on it.

```go
url, _ := c.OAuth(
	CLIENT_ID, CLIENT_SECRET, REDIRECT_URL, "profile",
)
```

After the user goes to `url` and grants your application permissions, you need to figure out a way for the user to input the second argument of the url to which they are redirected (ie: `REDIRECT_URL/?state=go-uber&code=AUTH_CODE`). You then need to

```go
c.SetAccessToken(AUTH_CODE)
```

At which point, feel free to

```go
profile, err := c.GetUserProfile()
if err != nil {
	fmt.Println(err)
} else {
	fmt.Println(profile)
}
```
