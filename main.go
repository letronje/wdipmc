package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/letronje/wdipmc/svy21"
)

type Carpark struct {
	gorm.Model
	Number        string `gorm:"primary_key;unique_index"`
	Address       string
	Latitude      float64
	Longitude     float64
	TotalLots     int64
	AvailableLots int64
}

const port = 80

func main() {
	db, err := gorm.Open("sqlite3", "carparks.db")
	// //db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		panic("failed to connect to database " + err.Error())
	}

	// Migrate the schema
	db.AutoMigrate(&Carpark{})

	carparkRows := getCarparks()

	for i, row := range carparkRows {
		latitude, longitude := svy21.ToLatLon(row.YCoord, row.XCoord)
		carpark := Carpark{
			Number:    row.Number,
			Address:   row.Address,
			Latitude:  latitude,
			Longitude: longitude,
		}
		db.FirstOrCreate(&carpark, "number = ?", row.Number)
		if i >= 3 {
			break
		}
	}

	r := gin.Default()

	// Get user value
	r.GET("/latlon", func(c *gin.Context) {
		north, _ := c.GetQuery("north")
		northf, _ := strconv.ParseFloat(north, 64)

		east, _ := c.GetQuery("east")
		eastf, _ := strconv.ParseFloat(east, 64)

		lat, lon := svy21.ToLatLon(northf, eastf)
		fmt.Println(lat, lon)

		c.JSON(
			http.StatusOK, gin.H{
				"latlon": fmt.Sprintf("%f,%f", lat, lon),
			})
	})

	r.Run(fmt.Sprintf(":%d", port))
}

type CarparkCSVRow struct {
	Number  string
	Address string
	XCoord  float64
	YCoord  float64
}

func getCarparks() []CarparkCSVRow {
	file, err := os.Open("hdb-carpark-information.csv")
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
	// 	"car_park_type",
	// 	"type_of_parking_system",
	// 	"short_term_parking",
	// 	"free_parking",
	// 	"night_parking",
	// 	"car_park_decks",
	// 	"gantry_height",
	// 	"car_park_basement",
	// }

	rows := []CarparkCSVRow{}
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

		rows = append(rows, CarparkCSVRow{
			Number:  record[0],
			Address: record[1],
			XCoord:  x,
			YCoord:  y,
		})
	}
	return rows
}
