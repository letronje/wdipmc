package carparkscsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Row struct {
	Number  string
	Address string
	XCoord  float64
	YCoord  float64
}

func Parse(filePath string) ([]Row, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []Row{}, fmt.Errorf("error opening carparks csv " + filePath)
	}

	reader := csv.NewReader(file)
	rows := []Row{}
	for rowNum := 1; ; rowNum++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []Row{}, err
		}

		// skip csv header
		if rowNum == 1 {
			continue
		}

		xCoord := strings.TrimSpace(record[2])
		x, err := strconv.ParseFloat(xCoord, 64)
		if err != nil {
			return []Row{}, fmt.Errorf("error parsing `x_coord` '%s' @ (row %d, col %d): %v", xCoord, rowNum, 3, err)
		}

		yCoord := strings.TrimSpace(record[3])
		y, err := strconv.ParseFloat(yCoord, 64)
		if err != nil {
			return []Row{}, fmt.Errorf("error parsing `y_coord` '%s' @ (row %d, col %d): %v", yCoord, rowNum, 4, err)
		}

		number := strings.TrimSpace(record[0])
		if number == "" {
			return []Row{}, fmt.Errorf("missing `car_park_no` @ (row %d, col %d)", rowNum, 1)
		}

		address := strings.TrimSpace(record[1])
		if address == "" {
			return []Row{}, fmt.Errorf("missing `address` @ (row %d, col %d)", rowNum, 2)
		}

		rows = append(rows, Row{
			Number:  number,
			Address: address,
			XCoord:  x,
			YCoord:  y,
		})
	}

	return rows, nil
}
