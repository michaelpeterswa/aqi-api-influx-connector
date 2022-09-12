package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AQIResponse struct {
	PrimaryPollutant string `json:"primary_pollutant"`
	Level            string `json:"level"`
	AQI              int64  `json:"aqi"`
}

func (r *RequestClient) GetCurrentAQI(endpoint string) (*AQIResponse, error) {
	url := fmt.Sprintf("%s/api/aqi", endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var aqi AQIResponse
	err = json.NewDecoder(resp.Body).Decode(&aqi)
	if err != nil {
		return nil, err
	}
	return &aqi, nil
}
