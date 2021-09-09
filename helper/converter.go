package helper

import (
	"fmt"
	"strings"
	"xrate/common/crate"
	"xrate/common/crate/provider"
	"xrate/config"
)

//CurrentConverterClient returns current converter client from os.environment
func CurrentConverterClient() (provider.Client, error) {
	client, err := crate.Open(config.GetClient().ApiProvider)
	if err != nil {
		return nil, err
	}

	return client, nil
}

//Conversion to hold convertion data
type Conversion struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Rate   float64 `json:"rate"`
	Value  float64 `json:"value"`
	Result float64 `json:"result"`
}

//Convert convert from curr to curr by given rates
func Convert(rates provider.Rate, from, to string, value float64) (Conversion, error) {
	var c Conversion

	if value < 1 {
		return c, fmt.Errorf("value can not less than 1")
	}

	//convert to uppercase
	from, to = strings.ToUpper(from), strings.ToUpper(to)
	fromBase, toBase := from == rates.Base, to == rates.Base
	if !(fromBase || toBase) {
		return c, fmt.Errorf("currencies not supported, %s to %s", from, to)
	}

	if from == to {
		return c, fmt.Errorf("cannot convert to same currency")
	}

	var crate float64
	if fromBase {
		v, ok := rates.Rates[to]
		if !ok {
			return c, fmt.Errorf("currency not supported: %s", from)
		}
		crate = v
	}

	if toBase {
		v, ok := rates.Rates[from]
		if !ok {
			return c, fmt.Errorf("currency not supported: %s", to)
		}
		crate = 1 / v
	}

	result := crate * value
	c = Conversion{
		From:   from,
		To:     to,
		Rate:   crate,
		Value:  value,
		Result: result,
	}

	return c, nil
}
