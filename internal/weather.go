package weather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ApiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func LoadApiKey(filename string) (ApiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ApiConfigData{}, err
	}

	var c ApiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return ApiConfigData{}, err
	}
	return c, nil
}

func Query(city string) (WeatherData, error) {
	apiconfig, err := LoadApiKey("config/.ApiConfig")
	if err != nil {
		return WeatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiconfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()
	var w WeatherData
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return WeatherData{}, err
	}
	return w, nil
}
