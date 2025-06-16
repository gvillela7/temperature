package services

import (
	"encoding/json"
	"errors"
	config "github.com/gvillela7/temperature/configs"
	"github.com/gvillela7/temperature/internal/data/response"
	"io"
	"net/http"
	"strings"
)

type TemperatureCelsius struct {
	Temp string `json:"temp"`
}

type Temperature struct {
	State string  `json:"state"`
	TempC float32 `json:"temp_c"`
	TempF float32 `json:"temp_f"`
	TempK float32 `json:"temp_k"`
}

type ViaCep struct {
	Estado string `json:"estado"`
	UF     string `json:"uf"`
	Erro   string `json:"erro,omitempty"`
}
type WeatherResponse struct {
	Current Current `json:"current"`
}
type Current struct {
	TempC float32 `json:"temp_c"`
}

func NewTemperature() *Temperature {
	return &Temperature{
		State: "",
		TempC: 0.0,
		TempF: 0.0,
		TempK: 0.0,
	}
}

func (t *Temperature) Celsius(cep string, w http.ResponseWriter) (*Temperature, error) {
	cfg := config.GetWeatherAPI()
	req, err := http.NewRequest(http.MethodGet, "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "error creating request for viacep.", nil)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		response.HttpResponse(w, http.StatusBadRequest, "request error.", nil)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "error read response viacep", nil)
		return nil, err
	}

	var viacep ViaCep
	if err := json.Unmarshal(body, &viacep); err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "error Unmarchal json viacep", nil)
		return nil, err
	}
	if viacep.Erro == "true" {
		response.HttpResponse(w, http.StatusNotFound, "zipcode not found.", nil)
		return nil, errors.New("zipcode not found")
	}
	if viacep.UF == "SP" {
		viacep.Estado = "Sao_Paulo"
	}

	state := strings.ReplaceAll(viacep.Estado, " ", "+")

	reqWeather, err := http.NewRequest(http.MethodGet, "http://api.weatherapi.com/v1/current.json?key="+cfg.Key+"&q="+state+"&aqi=no", nil)
	if err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "erro ao criar requisição para o weatherapi", nil)
		return nil, err
	}
	resWeather, err := http.DefaultClient.Do(reqWeather)
	if err != nil {
		response.HttpResponse(w, http.StatusNotFound, "não foi possível encontrar informações meteorológicas.", nil)
		return nil, err
	}
	defer resWeather.Body.Close()

	bodyWeather, err := io.ReadAll(resWeather.Body)
	if err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "erro ao ler resposta do weatherapi", nil)
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(bodyWeather, &weather); err != nil {
		response.HttpResponse(w, http.StatusInternalServerError, "erro ao decodificar resposta do weatherapi", nil)
		return nil, err
	}

	//tempC, _ := strconv.ParseFloat(weather.Current.TempC, 32)
	t.State = state
	t.TempC = weather.Current.TempC
	t.TempF, _ = t.Fahrenheit(t.TempC)
	t.TempK, _ = t.Kelvin(t.TempC)

	return t, nil
}

func (t *Temperature) Fahrenheit(celsius float32) (float32, error) {
	return celsius*1.8 + 32, nil
}

func (t *Temperature) Kelvin(celsius float32) (float32, error) {
	return celsius + 273.15, nil
}
