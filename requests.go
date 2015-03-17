package uber

//
// structs representing the necessary data for generating requests to the various
// endpoints
//

type authReq struct {
	// cannot be pointer because of reflection in `generateRequestURLHelper`
	auth
	responseType string `query:"response_type,required"`
	scope        string `query:"scope"`
	state        string `query:"state"`
}

type accReq struct {
	auth
	clientSecret string `query:"client_secret,required"`
	grantType    string `query:"grant_type,required"`
	code         string `query:"code,required"`
}

type requestReq struct {
	productID           string  `query:"product_id,required"`
	startLatitude       float64 `query:"start_latitude,required"`
	startLongitude      float64 `query:"start_longitude,required"`
	endLatitude         float64 `query:"end_latitude,required"`
	endLongitude        float64 `query:"end_longitude,required"`
	surgeConfirmationID string  `query:"surge_confirmation_id"`
}

type requestResp struct {
	Request
}

type requestMapResp struct {
	RequestID string `json:"request_id"`
	HRef      string `json:"href"`
}

type productsReq struct {
	latitude  float64 `query:"latitude,required"`
	longitude float64 `query:"longitude,required"`
}

// productsResp is the type that is returned from the `ProductEndpoint`
// This data definition is needed so that unmarshalling can actually happen.
type productsResp struct {
	Products []*Product `json:"products"`
}

type pricesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	endLatitude    float64 `query:"end_latitude,required"`
	endLongitude   float64 `query:"end_longitude,required"`
}

// pricesResp is the type that is returned from the `PriceEndpoint`
// This data definition is needed so that unmarshalling can actually happen.
type pricesResp struct {
	Prices []*Price `json:"prices"`
}

type timesReq struct {
	startLatitude  float64 `query:"start_latitude,required"`
	startLongitude float64 `query:"start_longitude,required"`
	customerUuid   string  `query:"customer_uuid"`
	productID      string  `query:"product_id"`
}

// timesResp is the type that is returned from the `PriceEndpoint`
// This data definition is needed so that unmarshalling can actually happen.
type timesResp struct {
	Times []*Time `json:"times"`
}

type historyReq struct {
	offset int `query:"offset,required"`
	limit  int `query:"limit,required"`
}
