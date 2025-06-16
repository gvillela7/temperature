package handler

import (
	"github.com/gvillela7/temperature/internal/data/response"
	"github.com/gvillela7/temperature/internal/services"
	"net/http"
	"strings"
)

func GetCep(w http.ResponseWriter, r *http.Request) {
	cepRequest := r.URL.Query().Get("cep")
	cep := strings.ReplaceAll(cepRequest, "-", "")
	if len(cep) != 8 {
		response.HttpResponse(w, http.StatusUnprocessableEntity, "invalid zipcode.", nil)
		return
	}

	service := services.NewTemperature()
	temperature, err := service.Celsius(cep, w)

	if err == nil {
		response.HttpResponse(w, http.StatusOK, "success", temperature)
		return
	}
}
