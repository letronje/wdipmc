package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/letronje/wdipmc/handlers"
	"github.com/letronje/wdipmc/stores/carparkstore"
)

const port = 80

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Println("Missing env var " + name)
		os.Exit(0)
	}
	return val
}

func main() {
	dsn := getEnvVar("CARPARK_DSN")

	var err error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := carparkstore.New(db)

	r := gin.Default()

	r.GET("/carparks/nearest", handlers.NearestCarParksHandler(store))

	r.Run(fmt.Sprintf(":%d", port))
}
