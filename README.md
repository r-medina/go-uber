go-uber
=======

Golang bindings for [Uber API](https://developer.uber.com/v1/endpoints/)

[![Build Status](https://travis-ci.org/r-medina/go-uber.png)](https://travis-ci.org/r-medina/go-uber)

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
  client := uber.NewClient(SERVER_TOKEN, ACCESS_TOKEN)
}
```

## Making Requests

Currently, the Uber API offers support for requesting information about products (e.g available cars), price estimates, time estimates, user ride history, and user info. All requests require a valid `server_token`. Requests that require latitude or longitude as arguments are float64's (and should be valid lat/lon's).

For more information about using this library, view the [Godocs](http://godoc.org/github.com/r-medina/go-uber).
