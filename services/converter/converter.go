package converter

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"xrate/common/crate/provider"
	"xrate/common/store/driver"
	"xrate/services"
	"xrate/services/scheduler"

	"github.com/pkg/errors"
)

// 1 hour
var delayTime = time.Second * 60 * 60

const (
	storekey = "rate"
)

//Conversion to hold convertion data
type Conversion struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Rate   float64 `json:"rate"`
	Value  float64 `json:"value"`
	Result float64 `json:"result"`
}

//IConverterService base converter methods
type IConverterService interface {
	services.Service
	//Rate returns currency rates
	Rate() (provider.Rate, error)
	//Convert convert currency
	Convert(from, to string, value float64) (Conversion, error)
}

type iconvService struct {
	sc     scheduler.ISchedulerService
	store  driver.Driver
	client provider.Client
}

//NewService returns new converter service instance
func NewService(sc scheduler.ISchedulerService, store driver.Driver, client provider.Client) IConverterService {
	return &iconvService{
		sc:     sc,
		store:  store,
		client: client,
	}
}

//Run implementing Services.Run
func (s *iconvService) Run(ctx context.Context) error {
	log.Println("Converter service started...")
	//get rates and update rate in background with scheduler service
	s.sc.AddTask(s.storeRate, delayTime)

	//block until context is done
	<-ctx.Done()
	return nil
}

//Rate implementing IConverterService.Rate
func (s *iconvService) Rate() (rate provider.Rate, err error) {
	i, err := s.store.Get(storekey)
	if err != nil {
		return rate, err
	}

	rate, ok := i.(provider.Rate)
	if !ok {
		return rate, fmt.Errorf("can not convert value to Rate instance")
	}

	return rate, nil
}

//Covert implementing IConverterService.Convert
func (s iconvService) Convert(from, to string, value float64) (Conversion, error) {
	var c Conversion

	if value < 1 {
		return c, fmt.Errorf("value can not less than 1")
	}

	//get rates
	rates, err := s.Rate()
	if err != nil {
		return c, errors.Wrap(err, "unable to get rates")
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
	if toBase {
		v, ok := rates.Rates[from]
		if !ok {
			return c, fmt.Errorf("currency not supported: %s", from)
		}
		crate = v
	}

	if fromBase {
		v, ok := rates.Rates[to]
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

func (s *iconvService) storeRate(ctx context.Context) error {
	rate, err := s.client.Rate(ctx)
	if err != nil {
		return err
	}
	return s.store.Set(storekey, rate)
}
