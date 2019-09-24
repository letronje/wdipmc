package carparkavailabilityapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type AvailabilityUpdate struct {
	Number        string
	TotalLots     int
	LotsAvailable int
}

type carparkInfo struct {
	TotalLots     string `json:"total_lots"`
	LotsAvailable string `json:"lots_available"`
}

type carparkData struct {
	Number      string        `json:"carpark_number"`
	CarparkInfo []carparkInfo `json:"carpark_info"`
}

type item struct {
	CarparkData []carparkData `json:"carpark_data"`
}

type availabilityResp struct {
	Items []item `json:"Items"`
}

const (
	url = "https://api.data.gov.sg/v1/transport/carpark-availability"
)

func Fetch() ([]AvailabilityUpdate, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []AvailabilityUpdate{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []AvailabilityUpdate{}, err
	}

	info := availabilityResp{}
	unmarshalErr := json.Unmarshal(bodyBytes, &info)
	if unmarshalErr != nil {
		return []AvailabilityUpdate{}, unmarshalErr
	}

	availabilities := []AvailabilityUpdate{}

	for _, carkparkData := range info.Items[0].CarparkData {
		available, _ := strconv.Atoi(carkparkData.CarparkInfo[0].LotsAvailable)
		total, _ := strconv.Atoi(carkparkData.CarparkInfo[0].TotalLots)
		availabilities = append(availabilities, AvailabilityUpdate{
			Number:        carkparkData.Number,
			LotsAvailable: available,
			TotalLots:     total,
		})
	}

	return availabilities, nil
}
