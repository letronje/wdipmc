package carparkstore

import (
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

type Store interface {
	Add(*Carpark) error
	UpdateAvailability(string, int, int) error
	NearestAvailable(latitude float64, longitude float64, page int, pagesize int) ([]Carpark, error)
}
