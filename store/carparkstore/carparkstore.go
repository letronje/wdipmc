package carparkstore

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Carpark struct {
	gorm.Model
	Number        string  `gorm:"primary_key;unique_index"`
	Address       string  `gorm:"not null"`
	Latitude      float64 `gorm:"not null"`
	Longitude     float64 `gorm:"not null"`
	TotalLots     int     `gorm:"not null"`
	AvailableLots int     `gorm:"not null;index:availableLots"`
}

var db *gorm.DB

func Init() {
	var err error
	//db, err = gorm.Open("sqlite3", "carparks.db")
	db, err = gorm.Open("mysql", "root:@/wdipmc?charset=utf8&parseTime=True&loc=Local")
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

func Add(carpark *Carpark) {
	db.FirstOrCreate(carpark, "number = ?", carpark.Number)
}

func UpdateAvailability(number string, lotsAvailable int, totalLots int) error {
	carpark := Carpark{}
	db.First(&carpark, "number = ?", number)
	if carpark == (Carpark{}) {
		return errors.New("Couldn't find carpark with number " + number)
	}
	carpark.TotalLots = totalLots
	carpark.AvailableLots = lotsAvailable
	db.Save(&carpark)
	return nil
}

func NearestAvailable(latitude float64, longitude float64, page int, pagesize int) []Carpark {
	orderBy := fmt.Sprintf("st_distance_sphere(point(%f, %f), point(longitude, latitude)) asc", longitude, latitude)
	carparks := []Carpark{}
	db.Where("available_lots > 0").Order(orderBy).Offset(pagesize * (page - 1)).Limit(pagesize).Find(&carparks)
	return carparks
}
