package services

import "net/http"

type TemperatureService interface {
	Fahrenheit(celsius float32) (float32, error)
	Celsius(cep string, w http.ResponseWriter) (*Temperature, error)
	Kelvin(celsius float32) (float32, error)
}
