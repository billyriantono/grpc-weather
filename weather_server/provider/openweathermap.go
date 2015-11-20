package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
)

const (
	openWeatherMapUrl = "http://api.openweathermap.org/data/2.5/weather"
)

type OpenWeatherMap struct {
	ApiKey string
}

func (p OpenWeatherMap) Query(q string) (WeatherInfo, error) {
	body, err := p.get(q)
	if err != nil {
		return EmptyResult, err
	}

	var result openWeatherMapResult
	json.Unmarshal(body, &result)

	return result.asWeatherInfo(), nil
}

func (p OpenWeatherMap) get(q string) ([]byte, error) {
	queryUrl := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapUrl, url.QueryEscape(q), p.ApiKey)

	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf("Unexpected response: %s", resp.Status)
		log.Println("Request failed:", respErr)
		return nil, respErr
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type openWeatherMapResult struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func (r openWeatherMapResult) asWeatherInfo() WeatherInfo {
	if r.found() {
		return WeatherInfo{
			Temperature: r.toCelcius(),
			Description: r.description(),
			Found:       true,
		}
	}
	return EmptyResult
}

func (r openWeatherMapResult) toCelcius() float64 {
	return math.Floor(r.Main.Kelvin-273.15) + 0.5
}

func (r openWeatherMapResult) description() string {
	return r.Weather[0].Description
}

func (r openWeatherMapResult) found() bool {
	return len(r.Weather) > 0
}
