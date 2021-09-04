package api

import (
	"net/http"
	"xrate/services/api/auth"
	"xrate/services/api/controller"
	"xrate/services/api/middleware"
	"xrate/services/converter"

	"github.com/gorilla/mux"
)

func newRouter(s *apiService) *mux.Router {
	r := mux.NewRouter()
	//add middleware
	addMiddlewares(s, r)

	//add routers
	addRateRouters(s.rate, r)
	addProjectRouters(s.auth, r)
	return r
}

func addRateRouters(s converter.IConverterService, r *mux.Router) {
	c := controller.NewRateHandler(s)
	r.HandleFunc("/convert", c.Convert).Methods(http.MethodGet)
}

func addProjectRouters(auth auth.Auth, r *mux.Router) {
	c := controller.NewAuthHandler(auth)
	r.HandleFunc("/project", c.Create).Methods(http.MethodGet)
}

func addMiddlewares(s *apiService, r *mux.Router) {
	r.Use(
		middleware.NewMAuth(s.auth).Auth,
	)
}
