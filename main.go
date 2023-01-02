package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	APIKey string `json:"apiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadAPIKey(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var config apiConfigData

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return apiConfigData{}, err
	}

	return config, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadAPIKey(".config")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.APIKey)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var wd weatherData

	if err := json.NewDecoder(resp.Body).Decode(&wd); err != nil {
		return weatherData{}, err
	}

	return wd, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}
