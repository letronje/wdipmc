package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/letronje/wdipmc/imports/carparkavailability"
)

func main() {
	dsn := getEnvVar("CARPARK_DSN")

	var err error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	success, failures, err := carparkavailability.Update(db)
	if err != nil {
		fmt.Printf("Error updating carpark availability, err: %v \n\n", err)
		return
	}

	fmt.Printf("Success: %d , Failures: %d\n", success, failures)
}

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Println("Missing env var " + name)
		os.Exit(0)
	}
	return val
}
