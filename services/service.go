package services

import (
	"context"

	"golang.org/x/sync/errgroup"
)

//Service base service methods
type Service interface {
	//Run execute Run instance service method
	Run(ctx context.Context) error
}

//RunServices execute each service and returns last error status
func RunServices(ctx context.Context, services ...Service) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		service := service
		group.Go(func() error {
			return service.Run(ctx)
		})
	}

	return group.Wait()
}
