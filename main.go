package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

const port = 80

func main() {
	r := gin.Default()

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.Run(fmt.Sprintf(":%d", port))
}
