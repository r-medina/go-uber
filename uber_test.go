package uber

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testClient      *Client
	testServerToken = "some_token"
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
)

func TestNewClient(t *testing.T) {
	testClient = NewClient(testServerToken, "" /* Empty access_token */)
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
