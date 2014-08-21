package uber

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

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	PromoCode string `json:"promo_code"`
}
