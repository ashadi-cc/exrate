package api

import (
	"context"
	"log"
	"net/http"
	"time"
	"xrate/config"

	"github.com/pkg/errors"
)

func runServer(ctx context.Context, cfg config.Server, service *apiService) error {
	r := addRouters(service)
	srv := &http.Server{
		Addr: cfg.Address,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IddleTimeout,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	ch := make(chan error, 1)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			ch <- errors.Wrap(err, "gracefully shutdown server error")
		}
		close(ch)
	}()

	log.Println("Server Listening on", cfg.Address)
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errors.Wrap(err, "error starting server")
	}

	return <-ch
}
