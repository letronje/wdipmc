package carparkstore

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type storeImpl struct {
	db *gorm.DB
}

func (s *storeImpl) init(db *gorm.DB) error {
	s.db = db
	s.db.AutoMigrate(&Carpark{})
	return nil
}

func (s *storeImpl) close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) Add(carpark *Carpark) error {
	s.db.FirstOrCreate(carpark, "number = ?", carpark.Number)
	return nil
}

func (s *storeImpl) UpdateAvailability(number string, lotsAvailable int, totalLots int) error {
	carpark := Carpark{}
	s.db.First(&carpark, "number = ?", number)
	if carpark == (Carpark{}) {
		return errors.New("Couldn't find carpark with number " + number)
	}
	carpark.TotalLots = totalLots
	carpark.AvailableLots = lotsAvailable
	s.db.Save(&carpark)
	return nil
}

func (s *storeImpl) NearestAvailable(latitude float64, longitude float64, page int, pagesize int) ([]Carpark, error) {
	orderBy := fmt.Sprintf("st_distance_sphere(point(%f, %f), point(longitude, latitude)) asc", longitude, latitude)
	carparks := []Carpark{}
	s.db.Where("available_lots > 0").Order(orderBy).Offset(pagesize * (page - 1)).Limit(pagesize).Find(&carparks)
	return carparks, nil
}

func New(db *gorm.DB) Store {
	impl := storeImpl{}
	_ = impl.init(db)
	return &impl
}
