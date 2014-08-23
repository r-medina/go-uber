package uber

type Client struct {
	server_token string
	access_token string
}

type Product struct {
	ProductId   string `json:"product_id"`
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	Capacity    int    `json:"capacity"`
	Image       string `json:"image"`
}

type Price struct {
	ProductId       string  `json:"product_id"`
	CurrencyCode    string  `json:"currency_code"`
	DisplayName     string  `json:"display_name"`
	Estimate        string  `json:"estimate"`
	LowEstimate     int     `json:"low_estimate"`
	HighEstimate    int     `json:"high_estimate"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
}

type Time struct {
	ProductId   string `json:"product_id"`
	DisplayName string `json:"display_name"`
	Estimate    int    `json:"estimate"`
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

type Location struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UserHistory struct {
	Offset  int     `json:"offset"`
	Limit   int     `json:"limit"`
	Count   int     `json:"count"`
	History []*Trip `json:"history"`
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
