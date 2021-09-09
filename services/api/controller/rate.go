package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"xrate/common/crate/provider"
	"xrate/helper"
)

type IRateService interface {
	GetRate() (provider.Rate, error)
}
type RateHandler struct {
	service IRateService
}

func NewRateHandler(service IRateService) *RateHandler {
	return &RateHandler{service: service}
}

func (h RateHandler) Convert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	value := r.URL.Query().Get("value")
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "value must be number"})
		return
	}

	rates, err := h.service.GetRate()
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	c, err := helper.Convert(rates, from, to, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(c)
}
