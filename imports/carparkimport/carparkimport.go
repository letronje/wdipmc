package carparkimport

import (
	"github.com/jinzhu/gorm"
	"github.com/letronje/wdipmc/csv/carparkscsv"
	"github.com/letronje/wdipmc/libs/svy21"
	"github.com/letronje/wdipmc/stores/carparkstore"
)

func Import(path string, db *gorm.DB) (int, error) {
	store := carparkstore.New(db)

	rows, err := carparkscsv.Parse(path)
	if err != nil {
		return 0, err
	}

	added := 0

	for _, row := range rows {
		latitude, longitude := svy21.ToLatLon(row.YCoord, row.XCoord)
		err = store.Add(&carparkstore.Carpark{
			Number:    row.Number,
			Address:   row.Address,
			Latitude:  latitude,
			Longitude: longitude,
		})
		if err != nil {
			return added, err
		}
		added++
	}

	return added, nil
}
