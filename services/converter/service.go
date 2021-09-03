package converter

import (
	"context"
	"fmt"
	"log"
	"time"
	"xrate/common/crate"
	"xrate/common/crate/provider"
	"xrate/common/store"
	"xrate/common/store/driver"
	"xrate/config"
	"xrate/services"
	"xrate/services/scheduler"

	//register driver
	_ "xrate/common/crate/client/fixer"
	_ "xrate/common/store/client/memory"
)

// 1 hour
var delayTime = time.Second * 60 * 60

const (
	storekey = "rate"
	currency = "USD"
)

type IConverterService interface {
	services.Service
	Rate() (float64, error)
}

type iconvService struct {
	sc     scheduler.ISchedulerService
	store  driver.Driver
	client provider.Client
}

func NewService(sc scheduler.ISchedulerService) IConverterService {
	client, err := crate.Open(config.GetClient().ApiProvider)
	if err != nil {
		panic(err)
	}

	store, err := store.Open(config.GetStore().Driver)
	if err != nil {
		panic(err)
	}

	return &iconvService{
		sc:     sc,
		store:  store,
		client: client,
	}
}

func (s *iconvService) Run(ctx context.Context) error {
	log.Println("Converter service started...")
	s.sc.AddTask(s.getRate, delayTime)
	return nil
}

func (s *iconvService) Rate() (float64, error) {
	rate, err := s.store.Get(storekey)
	if err != nil {
		return -1, err
	}

	v, ok := rate.(float64)
	if !ok {
		return -1, fmt.Errorf("can not convert value to float64")
	}

	return v, nil
}

func (s *iconvService) getRate(ctx context.Context) error {
	rate, err := s.client.Rate(ctx, currency)
	if err != nil {
		return err
	}
	log.Println("rate:", rate)
	return s.store.Set(storekey, rate)
}
