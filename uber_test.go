package uber

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testClient      *Client
	testServerToken = "some_token"
	testAccessToken = "bearer_token"
	testProducts    = []*Product{
		&Product{
			ProductId:   "1",
			Description: "The Original Uber",
			DisplayName: "UberBLACK",
			Capacity:    4,
			Image:       "http://...",
		},
	}
	testPrices = []*Price{
		&Price{
			ProductId:       "1",
			CurrencyCode:    "USD",
			DisplayName:     "UberBlack",
			Estimate:        "$23-29",
			LowEstimate:     23,
			HighEstimate:    29,
			SurgeMultiplier: 1.25,
		},
	}
	testTimes = []*Time{
		&Time{
			ProductId:   "1",
			DisplayName: "UberBLACK",
			Estimate:    400,
		},
	}
	testUserActivity = &UserActivity{
		Offset: 0,
		Limit:  2,
		Count:  1,
		History: []*Trip{
			&Trip{
				Uuid:        "7354db54-cc9b-4961-81f2-0094b8e2d215",
				RequestTime: 1401884467,
				ProductId:   "edf5e5eb-6ae6-44af-bec6-5bdcf1e3ed2c",
				Status:      "completed",
				Distance:    0.0279562,
				StartTime:   1401884646,
				StartLocation: &Location{
					Address:   "706 Mission St, San Francisco, CA",
					Latitude:  37.7860099,
					Longitude: -122.4025387,
				},
				EndTime: 1401884732,
				EndLocation: &Location{
					Address:   "1455 Market Street, San Francisco, CA",
					Latitude:  37.7758179,
					Longitude: -122.4180285,
				},
			},
		},
	}
	testUserProfile = &User{
		FirstName: "Uber",
		LastName:  "Developer",
		Email:     "developer@uber.com",
		Picture:   "https://...",
		PromoCode: "teypo",
	}
)

func TestNewClient(t *testing.T) {
	testClient = NewClient(testServerToken, testAccessToken)
	if testClient.serverToken != testServerToken {
		t.Fatal(fmt.Sprintf("Client.serverToken %s does not match %s", testClient.serverToken, testServerToken))
	}
}

func TestGetProducts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getProductsHandler))
	defer server.Close()
	UBER_API_ENDPOINT = server.URL

	_, err := testClient.GetProducts(123.0, 456.0)
	if err != nil {
		t.Fatal(err)
	}
}

func getProductsHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := json.Marshal(testProducts)
	rw.Write(body)
}

func TestGetPrices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getPricesHandler))
	defer server.Close()
	UBER_API_ENDPOINT = server.URL

	_, err := testClient.GetPrices(123.0, 456.0, 234.0, 567.0)
	if err != nil {
		t.Fatal(err)
	}
}

func getPricesHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := json.Marshal(testPrices)
	rw.Write(body)
}

func TestGetTimes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getTimesHandler))
	defer server.Close()
	UBER_API_ENDPOINT = server.URL

	_, err := testClient.GetTimes(123.0, 456.0, "" /* uuid */, "" /* productId */)
	if err != nil {
		t.Fatal(err)
	}
}

func getTimesHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := json.Marshal(testTimes)
	rw.Write(body)
}

func TestGetUserActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getUserActivityHandler))
	defer server.Close()
	UBER_API_ENDPOINT = server.URL

	_, err := testClient.GetUserActivity(0 /* offset */, 2 /* count */)
	if err != nil {
		t.Fatal(err)
	}
}

func getUserActivityHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := json.Marshal(testUserActivity)
	rw.Write(body)
}

func TestGetUserProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getUserProfileHandler))
	defer server.Close()
	UBER_API_ENDPOINT = server.URL

	_, err := testClient.GetUserProfile()
	if err != nil {
		t.Fatal(err)
	}
}

func getUserProfileHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := json.Marshal(testUserProfile)
	rw.Write(body)
}

func TestGet(t *testing.T) {
	t.Fatal("no test")
}

func TestSendRequestWithAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(sendRequestWithAuthorizationHandler))
    defer server.Close()

    // Send with only serverToken i.e. oauth is false
    res, err := testClient.sendRequestWithAuthorization(server.URL, false);
    if err != nil {
        t.Fatal(err)
    }
    auth := res.Request.Header.Get("Authorization")
    if auth == "" || auth != fmt.Sprintf("Token %s", testServerToken) {
        t.Fatal("Server token not found in header")
    }

    // Send with only accessToken i.e. oauth is true
    res, err = testClient.sendRequestWithAuthorization(server.URL, true);
    if err != nil {
        t.Fatal(err)
    }
    auth = res.Request.Header.Get("Authorization")
    if auth == "" || auth != fmt.Sprintf("Bearer %s", testAccessToken) {
        t.Fatal("Access token not found in header")
    }
}

func sendRequestWithAuthorizationHandler(rw http.ResponseWriter, req *http.Request) {
    rw.Write([]byte{0})
}

func TestGenerateRequestUrl(t *testing.T) {
	t.Fatal("no test")
}

func TestGenerateRequestUrlHelper(t *testing.T) {
	t.Fatal("no test")
}

// TODO: test `get`, `generateRequestUrl`, and `generateRequestUrlHelper`
