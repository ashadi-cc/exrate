package api

import (
	"context"
	"fmt"
	"log"
	"time"
	"xrate/common/crate/provider"
	"xrate/common/store/driver"
	"xrate/config"
	"xrate/services"
	"xrate/services/api/auth"
	"xrate/services/scheduler"
)

const rateStoreKey = "rate"

// 1 hour
var delayTime = time.Second * 30

//IApiService base methods api service interface
type IApiService interface {
	services.Service
	GetRate() (rate provider.Rate, err error)
}

type apiService struct {
	store     driver.Driver
	auth      auth.Auth
	scheduler scheduler.ISchedulerService
	client    provider.Client
}

//returns new instance api service
func NewService(sceduler scheduler.ISchedulerService, client provider.Client, auth auth.Auth, store driver.Driver) IApiService {
	return &apiService{
		scheduler: sceduler,
		store:     store,
		auth:      auth,
		client:    client,
	}
}

//Run implementing services.Service
func (s *apiService) Run(ctx context.Context) error {
	//get rates in backgrond and store it
	s.getRates()

	//run api service
	log.Println("API service started...")
	return runServer(ctx, config.GetServer(), s)
}

func (s *apiService) getRates() {
	s.scheduler.AddTask(s.storeRate, delayTime)
}

func (s *apiService) storeRate(ctx context.Context) error {
	rate, err := s.client.Rate(ctx)
	if err != nil {
		return err
	}
	return s.store.Set(rateStoreKey, rate)
}

func (s *apiService) GetRate() (rate provider.Rate, err error) {
	i, err := s.store.Get(rateStoreKey)
	if err != nil {
		return rate, err
	}

	rate, ok := i.(provider.Rate)
	if !ok {
		return rate, fmt.Errorf("can not convert value to Rate instance")
	}

	return rate, nil
}
