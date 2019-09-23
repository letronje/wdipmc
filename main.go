package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/letronje/wdipmc/svy21"
)

const port = 80

func main() {
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
