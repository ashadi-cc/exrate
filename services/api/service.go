package api

import (
	"context"
	"log"
	"xrate/config"
	"xrate/services"
	"xrate/services/converter"
)

//IApiService base methods api service interface
type IApiService interface {
	services.Service
}

type apiService struct {
	rate converter.IConverterService
}

//returns new instance api service
func NewService(rate converter.IConverterService) IApiService {
	return &apiService{
		rate: rate,
	}
}

//Run implementing services.Service
func (s *apiService) Run(ctx context.Context) error {
	log.Println("API service started...")

	return runServer(ctx, config.GetServer())
}
