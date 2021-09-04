package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"xrate/helper"
	"xrate/services"
	"xrate/services/api"
	"xrate/services/converter"
	"xrate/services/scheduler"
)

//NewApi returns new instance api app
func NewApi() App {
	app := NewApp(context.Background())

	app.SetactionFunc(func(ctx context.Context) error {
		return runApi(ctx)
	})
	return app
}

//runApi run api handler
func runApi(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("graceful shutdowns...")
		cancel()
	}()

	client, err := helper.CurrentConverterClient()
	if err != nil {
		return err
	}

	store, err := helper.CurrentStore()
	if err != nil {
		return err
	}

	schedulerService := scheduler.NewService()
	converterService := converter.NewService(schedulerService, store, client)
	apiService := api.NewService(converterService)

	//register and run services
	return services.RunServices(ctx, apiService, schedulerService, converterService)
}
