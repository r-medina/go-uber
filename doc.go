// This package provides an API client for the Uber API. It exposes methods to get
// information about Uber products, estimates, times, and users as well as actually
// requesting an Uber.
//
// A lot of documentation will be pulled directly from
// https://developer.uber.com/v1/endpoints.
//
// Organization
//
// This package is organized into a few files.
//
// 1. `uber.go` contains all the exported types (that directly reflect some json object
// the Uber API returns) that this package contains. This file also has global constants
// and variables. Finally, `uber.go` contains a few error types that are used in the
// package itself.
//
// 2. `client.go` contains the definition of `Client` (the type with which the user
// interacts). Aside from the constructor for the client, this file contains low-level
// functions for generating and executing HTTP requests to the Uber API.
//
// 3. `endpoints.go` contains the definitions of the exported methods on `Client` that
// call the Uber API endpoints. This is the meat of this package's API.
//
// 4. `auth.go` contains all the functions related to authorizing your app.
//
// 5. `requests.go` contains a plethora of unexported types needed to make requests and
// parse responses.
//
// TODO
//
// Write tests.
package uber
