package api

import (
	"net/http"
	"xrate/services/api/controller"

	"github.com/gorilla/mux"
)

func addRouters(s *apiService) *mux.Router {
	r := mux.NewRouter()
	addRateRouters(s, r)
	return r
}

func addRateRouters(s *apiService, r *mux.Router) {
	c := controller.NewRateHandler(s.rate)
	r.HandleFunc("/rate", c.Convert).Methods(http.MethodGet)
}
