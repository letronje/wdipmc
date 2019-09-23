package carparkstore

import (
	"github.com/jinzhu/gorm"
)

type Carpark struct {
	gorm.Model
	Number        string `gorm:"primary_key;unique_index"`
	Address       string
	Latitude      float64
	Longitude     float64
	TotalLots     int
	AvailableLots int
}

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open("sqlite3", "carparks.db")
	if err != nil {
		panic("failed to connect to database " + err.Error())
	}
	db.AutoMigrate(&Carpark{})
}

func Close() {
	err := db.Close()
	if err != nil {
		panic("failed to close database" + err.Error())
	}
}

func AddCarpark(carpark *Carpark) {
	db.FirstOrCreate(carpark, "number = ?", carpark.Number)
}

func UpdateCarparkAvailability(number string, lotsAvailable int, totalLots int) {
	carpark := Carpark{}
	db.First(&carpark, "number = ?", number)
	carpark.TotalLots = totalLots
	carpark.AvailableLots = lotsAvailable
	db.Save(&carpark)
}
