package api

import (
	"context"
	"log"
	"xrate/common/store/driver"
	"xrate/config"
	"xrate/services"
	"xrate/services/api/auth"
	"xrate/services/converter"
)

//IApiService base methods api service interface
type IApiService interface {
	services.Service
}

type apiService struct {
	store driver.Driver
	rate  converter.IConverterService
	auth  auth.Auth
}

//returns new instance api service
func NewService(rate converter.IConverterService, auth auth.Auth, store driver.Driver) IApiService {
	return &apiService{
		rate:  rate,
		store: store,
		auth:  auth,
	}
}

//Run implementing services.Service
func (s *apiService) Run(ctx context.Context) error {
	log.Println("API service started...")

	return runServer(ctx, config.GetServer(), s)
}
