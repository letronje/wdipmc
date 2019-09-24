package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/letronje/wdipmc/imports/carparkimport"
)

func main() {
	dsn := getEnvVar("CARPARK_DSN")
	csvPath := getEnvVar("CARPARK_CSV")

	var err error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	imported, err := carparkimport.Import(csvPath, db)
	if err != nil {
		fmt.Printf("Error importing carparks from '%s', err: %v \n", csvPath, err)
		return
	}

	fmt.Printf("Carparks after import: %d", imported)
}

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Println("Missing env var " + name)
		os.Exit(0)
	}
	return val
}
