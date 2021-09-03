package api

import (
	"context"
	"net/http"
	"time"
	"xrate/config"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func runServer(ctx context.Context, cfg config.API) error {
	r := mux.NewRouter()

	port := "8080"
	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
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

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errors.Wrap(err, "error starting server")
	}

	return <-ch

}
