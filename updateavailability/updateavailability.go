package updateavailability

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/letronje/wdipmc/store/carparkstore"
)

type CarparkUpdate struct {
	Number        string
	TotalLots     int
	LotsAvailable int
}

func Update() error {
	updates, err := getUpdates()
	if err != nil {
		return err
	}

	carparkstore.Init()
	defer carparkstore.Close()

	spew.Dump(updates)
	for _, update := range updates {
		carparkstore.UpdateCarparkAvailability(update.Number, update.LotsAvailable, update.TotalLots)
	}

	return nil
}

func getUpdates() ([]CarparkUpdate, error) {
	resp, err := http.Get("https://api.data.gov.sg/v1/transport/carpark-availability")
	if err != nil {
		return []CarparkUpdate{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []CarparkUpdate{}, err
	}

	//fmt.Println(string(bodyBytes))

	type CarparkInfo struct {
		TotalLots     string `json:"total_lots"`
		LotsAvailable string `json:"lots_available"`
	}

	type CarparkData struct {
		Number      string        `json:"carpark_number"`
		CarparkInfo []CarparkInfo `json:"carpark_info"`
	}

	type Item struct {
		CarparkData []CarparkData `json:"carpark_data"`
	}

	type Updates struct {
		Items []Item `json:"items"`
	}

	updateInfo := Updates{}
	unmarshalErr := json.Unmarshal(bodyBytes, &updateInfo)
	if unmarshalErr != nil {
		return []CarparkUpdate{}, unmarshalErr
	}

	fmt.Println(len(updateInfo.Items))

	updates := []CarparkUpdate{}
	for _, carkparkData := range updateInfo.Items[0].CarparkData {
		available, _ := strconv.Atoi(carkparkData.CarparkInfo[0].LotsAvailable)
		total, _ := strconv.Atoi(carkparkData.CarparkInfo[0].TotalLots)
		updates = append(updates, CarparkUpdate{
			Number:        carkparkData.Number,
			LotsAvailable: available,
			TotalLots:     total,
		})
	}

	return updates, nil
}
