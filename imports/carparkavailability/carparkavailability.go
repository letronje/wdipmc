package carparkavailability

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/letronje/wdipmc/ext/carparkavailabilityapi"
	"github.com/letronje/wdipmc/stores/carparkstore"
)

func Update(db *gorm.DB) (int, int, error) {
	store := carparkstore.New(db)

	updates, err := carparkavailabilityapi.Fetch()
	if err != nil {
		return 0, 0, err
	}

	success := 0
	failures := 0

	for _, update := range updates {
		err := store.UpdateAvailability(update.Number, update.LotsAvailable, update.TotalLots)
		if err != nil {
			fmt.Printf("Error updating availability for carpark '%s', err: %v\n", update.Number, err)
			failures++
			continue
		}
		success++
	}

	return success, failures, nil
}
