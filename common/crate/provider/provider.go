package provider

import "context"

//Client to hold 3rd api methods
type Client interface {
	//Rate returns EUR rate by given currency symbol
	Rate(ctx context.Context, curr string) (float64, error)
}
