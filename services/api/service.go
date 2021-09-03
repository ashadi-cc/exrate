package api

import (
	"context"
	"log"
	"xrate/config"
	"xrate/services"
	"xrate/services/scheduler"
)

//IApiService base methods api service interface
type IApiService interface {
	services.Service
}

type apiService struct {
	sc scheduler.ISchedulerService
}

//returns new instance api service
func NewService(sc scheduler.ISchedulerService) IApiService {
	return &apiService{
		sc: sc,
	}
}

//Run implementing services.Service
func (s *apiService) Run(ctx context.Context) error {
	log.Println("API service started...")

	return runServer(ctx, config.API{})
}
