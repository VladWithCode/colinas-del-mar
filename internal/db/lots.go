package db

import (
	"context"
	"time"
)

type Lot struct {
	Id        string             `json:"id" db:"id"`
	Lt        string             `json:"lt" db:"lt"`
	Mz        string             `json:"mz" db:"mz"`
	Available bool               `json:"available" db:"available"`
	Area      float64            `json:"area" db:"area"`
	Measures  map[string]float64 `json:"measures" db:"measures"`
	Type      string             `json:"type" db:"type"`
	Price     float64            `json:"price" db:"price"`
}

func CreateLot(lot *Lot) error {
	conn, err := GetPool()
	if err != nil {
		return err
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = conn.Exec(
		ctx,
		"INSERT INTO properties (id, lt, mz, available, area, measures, type, price) VALUES $1, $2, $3, $4, $5, $6, $7",
		lot.Id,
		lot.Lt,
		lot.Mz,
		lot.Area,
		lot.Measures,
		lot.Type,
		lot.Price,
	)

	return nil
}
