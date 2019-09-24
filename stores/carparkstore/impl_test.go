package carparkstore

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func getDB() *gorm.DB {
	dsn := getEnvVar("CARPARK_DSN_TEST")

	var err error
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func cleanDB(db *gorm.DB) {
	db.Exec("DELETE FROM carparks")
}

func Test_storeImpl_Add_HappyPath(t *testing.T) {
	db := getDB()
	defer db.Close()

	s := New(db)
	cleanDB(db)

	cp := Carpark{
		Number:        "123",
		Address:       "Home",
		Latitude:      1,
		Longitude:     1,
		TotalLots:     0,
		AvailableLots: 0,
	}

	_ = s.Add(&cp)

	carpark := Carpark{}
	db.First(&carpark, "number = ?", cp.Number)
	if carpark.Number != cp.Number {
		t.Errorf("can't find added carpark, inserted: %+v, found: %+v", cp, carpark)
	}
}

func Test_storeImpl_Add_Existing(t *testing.T) {
	db := getDB()
	defer db.Close()

	s := New(db)
	cleanDB(db)

	cp1 := Carpark{
		Number:        "123",
		Address:       "Home",
		Latitude:      1,
		Longitude:     1,
		TotalLots:     0,
		AvailableLots: 0,
	}

	cp2 := Carpark{
		Number:        "123",
		Address:       "Work",
		Latitude:      1,
		Longitude:     1,
		TotalLots:     0,
		AvailableLots: 0,
	}
	_ = s.Add(&cp1)
	_ = s.Add(&cp2)

	carpark := Carpark{}
	db.First(&carpark, "number = ?", cp1.Number)
	if carpark.Number != cp1.Number || carpark.Address != cp1.Address {
		t.Errorf("can't find added carpark, inserted: %+v, found: %+v", cp1, carpark)
	}
}

func Test_storeImpl_UpdateAvailability_HappyPath(t *testing.T) {
	db := getDB()
	defer db.Close()

	s := New(db)
	cleanDB(db)

	cp := Carpark{
		Number:        "123",
		Address:       "Home",
		Latitude:      1,
		Longitude:     1,
		TotalLots:     0,
		AvailableLots: 0,
	}

	db.Create(&cp)

	_ = s.UpdateAvailability(cp.Number, 10, 100)

	carpark := Carpark{}
	db.First(&carpark, "number = ?", cp.Number)
	if carpark.TotalLots != 100 {
		t.Errorf("total lots doesn't match, expected: %v, got: %v", 100, carpark.TotalLots)
	}
	if carpark.AvailableLots != 10 {
		t.Errorf("available lots doesn't match, expected: %v, got: %v", 10, carpark.AvailableLots)
	}
}

func Test_storeImpl_NearestAvailable_HappyPath(t *testing.T) {
	db := getDB()
	defer db.Close()

	s := New(db)
	cleanDB(db)

	carparks := []Carpark{
		{
			Number:        "BRM6",
			Address:       "BLK 116 JALAN TENTERAM",
			Latitude:      1.3277824471296569,
			Longitude:     103.86119305602381,
			TotalLots:     467,
			AvailableLots: 205,
		},
		{
			Number:        "KB20",
			Address:       "BLK 49A WHAMPOA SOUTH",
			Latitude:      1.3237488106749886,
			Longitude:     103.86644226080786,
			TotalLots:     223,
			AvailableLots: 146,
		},
		{
			Number:        "SG3",
			Address:       "BLK 20/23 ST GEORGE RD",
			Latitude:      1.3251463335586366,
			Longitude:     103.86217635268844,
			TotalLots:     141,
			AvailableLots: 30,
		},
	}

	for _, cp := range carparks {
		db.Create(&cp)
	}

	nearest, _ := s.NearestAvailable(1.3275, 103.8657, 1, 3)
	if nearest[0].Number != "KB20" || nearest[1].Number != "SG3" || nearest[2].Number != "BRM6" {
		t.Errorf("mismatch in nearest carparks, got: %+v", nearest)
	}
}

func Test_storeImpl_NearestAvailable_Pagination(t *testing.T) {
	db := getDB()
	defer db.Close()

	s := New(db)
	cleanDB(db)

	carparks := []Carpark{
		{
			Number:        "BRM6",
			Address:       "BLK 116 JALAN TENTERAM",
			Latitude:      1.3277824471296569,
			Longitude:     103.86119305602381,
			TotalLots:     467,
			AvailableLots: 205,
		},
		{
			Number:        "KB20",
			Address:       "BLK 49A WHAMPOA SOUTH",
			Latitude:      1.3237488106749886,
			Longitude:     103.86644226080786,
			TotalLots:     223,
			AvailableLots: 146,
		},
		{
			Number:        "PP9T",
			Address:       "BLK 145 POTONG PASIR AVENUES 2/3",
			Latitude:      1.3321940115902131,
			Longitude:     103.86605975349579,
			TotalLots:     131,
			AvailableLots: 87,
		},
		{
			Number:        "SG3",
			Address:       "BLK 20/23 ST GEORGE RD",
			Latitude:      1.3251463335586366,
			Longitude:     103.86217635268844,
			TotalLots:     141,
			AvailableLots: 30,
		},
	}

	for _, cp := range carparks {
		db.Create(&cp)
	}

	nearest, _ := s.NearestAvailable(1.3275, 103.8657, 2, 2)
	if nearest[0].Number != "BRM6" || nearest[1].Number != "PP9T" {
		t.Errorf("mismatch in nearest carparks, got: %+v", nearest)
	}
}

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Println("Missing env var " + name)
		os.Exit(0)
	}
	return val
}
