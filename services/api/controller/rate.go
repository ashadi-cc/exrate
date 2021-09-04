package controller

import (
	"encoding/json"
	"fmt"
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

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	c, err := convertVal(rate, from, to, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(c)
}

func convertVal(rate float64, from, to string, value float64) (convert, error) {
	result := convert{}
	fok := availableCurrency(from)
	tok := availableCurrency(to)
	if !fok || !tok {
		return result, fmt.Errorf("currency not supported")
	}

	if to == from {
		return result, fmt.Errorf("can't convert to same currency")
	}

	//simple converter eur <> usd with if condition
	var crate float64
	if from == "usd" && to == "eur" {
		crate = rate
	}
	if from == "eur" && to == "usd" {
		crate = 1 / rate
	}

	res := crate * value
	result.From = from
	result.To = to
	result.Rate = crate
	result.Value = res

	return result, nil
}

func availableCurrency(curr string) bool {
	curss := map[string]bool{"usd": true, "eur": true}
	return curss[curr]
}

type convert struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Rate  float64 `json:"rate"`
	Value float64 `json:"value"`
}
