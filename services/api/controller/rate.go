package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"xrate/services/converter"
)

type RateHandler struct {
	service converter.IConverterService
}

func NewRateHandler(service converter.IConverterService) *RateHandler {
	return &RateHandler{service: service}
}

func (h RateHandler) Convert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//get base rate
	rate, err := h.service.Rate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "can not load rate"})
		return
	}

	value := r.URL.Query().Get("value")
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "value must be number"})
		return

	}

	currency := r.URL.Query().Get("curr")
	c, err := convertVal(rate, currency, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(c)
}

func convertVal(rate float64, currency string, value float64) (convert, error) {
	result := convert{}

	return result, nil
}

type convert struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
	Value    float64 `json:"value"`
}
