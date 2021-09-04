package fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"xrate/common/crate/provider"
	"xrate/config"

	"github.com/pkg/errors"
)

const baseApiUrl = "http://data.fixer.io/api"

type client struct {
	cfg config.Client
}

type result struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

func (c *client) Rate(ctx context.Context) (provider.Rate, error) {
	var rate provider.Rate
	apiUrl := fmt.Sprintf("%s/latest?access_key=%s", baseApiUrl, c.cfg.ApiKey)
	log.Println("fixer: get rate from", apiUrl)
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return rate, errors.Wrapf(err, "unable to create request %s", apiUrl)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return rate, errors.Wrap(err, "unable to create request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return rate, errors.Wrap(err, "unable to read body")
	}

	if res.StatusCode != http.StatusOK {
		return rate, fmt.Errorf("unable to get rate %s", string(body))
	}

	var r result
	if err := json.Unmarshal(body, &r); err != nil {
		return rate, errors.Wrap(err, "unable to decode body")
	}

	if !r.Success {
		return rate, fmt.Errorf("%v", r)
	}

	rate = provider.Rate{
		Base:  r.Base,
		Rates: r.Rates,
	}

	_ = json.NewEncoder(log.Writer()).Encode(rate)
	return rate, nil
}
