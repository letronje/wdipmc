package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/letronje/wdipmc/store/carparkstore"
)

const port = 80

func main() {
	carparkstore.Init()
	defer carparkstore.Close()

	r := gin.Default()

	// Get user value
	r.GET("/carparks/nearest", func(c *gin.Context) {
		latitudeStr := c.Query("latitude")
		if latitudeStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"missing_field": "latitude",
			})
			return
		}
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"invalid_field": "latitude",
			})
			return
		}

		longitudeStr := c.Query("longitude")
		if longitudeStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"missing_field": "longitude",
			})
		}
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"invalid_field": "longitude",
			})
			return
		}

		pageStr := c.Query("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"invalid_field": "page",
			})
			return
		}

		pagesizeStr := c.Query("per_page")
		if pagesizeStr == "" {
			pagesizeStr = "10"
		}
		pagesize, err := strconv.Atoi(pagesizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"invalid_field": "per_page",
			})
			return
		}

		carparks := carparkstore.NearestAvailable(latitude, longitude, page, pagesize)

		response := []gin.H{}
		for _, carpark := range carparks {
			response = append(response, gin.H{
				"address":        carpark.Address,
				"latitude":       carpark.Latitude,
				"longitude":      carpark.Longitude,
				"total_lots":     carpark.TotalLots,
				"available_lots": carpark.AvailableLots,
			})
		}

		c.JSON(http.StatusOK, response)
	})

	r.Run(fmt.Sprintf(":%d", port))

	// carparkimport.Import("hdb-carpark-information.csv")
	// err := updateavailability.Update()
	// fmt.Println(err)
	//
	// carparkstore.Init()
	// defer carparkstore.Close()
	//
	// spew.Dump(carparkstore.NearestAvailable(1.3275, 103.8657)[1])
}
