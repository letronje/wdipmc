package carparkimport

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/letronje/wdipmc/store/carparkstore"
	"github.com/letronje/wdipmc/svy21"
)

type CSVRow struct {
	Number  string
	Address string
	XCoord  float64
	YCoord  float64
}

func Import(path string) {
	carparkstore.Init()
	defer carparkstore.Close()

	rows := carparkRows(path)

	for _, row := range rows {
		latitude, longitude := svy21.ToLatLon(row.YCoord, row.XCoord)
		carparkstore.Add(&carparkstore.Carpark{
			Number:    row.Number,
			Address:   row.Address,
			Latitude:  latitude,
			Longitude: longitude,
		})
	}

	fmt.Println(len(rows))
}

func carparkRows(path string) []CSVRow {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error open carparks csv")
		os.Exit(0)
	}

	r := csv.NewReader(file)
	// headers := []string{
	// 	"car_park_no",
	// 	"address",
	// 	"x_coord",
	// 	"y_coord",
	// 	...
	// 	...
	// 	...
	// }

	rows := []CSVRow{}
	for i := 0; ; i++ {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if i == 0 {
			continue
		}

		x, _ := strconv.ParseFloat(record[2], 64)
		y, _ := strconv.ParseFloat(record[3], 64)

		//pp.Println(record)

		rows = append(rows, CSVRow{
			Number:  record[0],
			Address: record[1],
			XCoord:  x,
			YCoord:  y,
		})
	}

	return rows
}
