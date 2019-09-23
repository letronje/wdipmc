package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/letronje/wdipmc/carparkimport"
	"github.com/letronje/wdipmc/updateavailability"
)

const port = 80

func main() {
	// r := gin.Default()
	//
	// // Get user value
	// r.GET("/latlon", func(c *gin.Context) {
	// 	north, _ := c.GetQuery("north")
	// 	northf, _ := strconv.ParseFloat(north, 64)
	//
	// 	east, _ := c.GetQuery("east")
	// 	eastf, _ := strconv.ParseFloat(east, 64)
	//
	// 	lat, lon := svy21.ToLatLon(northf, eastf)
	// 	fmt.Println(lat, lon)
	//
	// 	c.JSON(
	// 		http.StatusOK, gin.H{
	// 			"latlon": fmt.Sprintf("%f,%f", lat, lon),
	// 		})
	// })
	//
	// r.Run(fmt.Sprintf(":%d", port))

	carparkimport.Import("hdb-carpark-information.csv")
	err := updateavailability.Update()
	fmt.Println(err)
}
