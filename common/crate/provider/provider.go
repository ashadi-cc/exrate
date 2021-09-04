package provider

import "context"

//Client to hold 3rd api methods
type Client interface {
	//Rate returns EUR rate
	Rate(ctx context.Context) (Rate, error)
}

//Rate to hold currency rate
type Rate struct {
	//Base base currency
	Base string `json:"base"`
	//Rates list rate by currencies
	Rates map[string]float64 `json:"rates"`
}
