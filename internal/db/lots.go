package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Lot struct {
	Id          string  `json:"id" db:"id"`
	Lt          string  `json:"lt" db:"lt"`
	Mz          string  `json:"mz" db:"mz"`
	Available   bool    `json:"available" db:"available"`
	Area        float64 `json:"area" db:"area"`
	Measures    string  `json:"measures" db:"measures"`
	Type        string  `json:"type" db:"type"`
	PriceCash   float64 `json:"priceCash" db:"price_cash"`
	PriceCredit float64 `json:"priceCredit" db:"price_credit"`
	BatchId     string  `json:"batchId" db:"batch_id"`
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
		"INSERT INTO properties (id, lt, mz, available, area, measures, type, price_cash, price_credit) VALUES $1, $2, $3, $4, $5, $6, $7, $8",
		lot.Id,
		lot.Lt,
		lot.Mz,
		lot.Area,
		lot.Measures,
		lot.Type,
		lot.PriceCash,
	)

	return nil
}

func BatchCreateLot(b *pgx.Batch) error {
	conn, err := GetPool()
	if err != nil {
		return err
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	br := conn.SendBatch(ctx, b)
	defer br.Close()

	for i, l := 0, b.Len(); i < l; i++ {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func QueueCreateLot(b *pgx.Batch, lot *Lot) {
	b.Queue(
		"INSERT INTO properties (id, lt, mz, available, area, measures, type, price_cash, price_credit, batch_id) VALUES $1, $2, $3, $4, $5, $6, $7, $8, $9, $10",
		lot.Id,
		lot.Lt,
		lot.Mz,
		lot.Area,
		lot.Measures,
		lot.Type,
		lot.PriceCash,
		lot.PriceCredit,
		lot.BatchId,
	)
}

// Silently tries to remove lots when batch insert fails
func RemoveFailedBatchLots(id string) {
	conn, err := GetPool()
	if err != nil {
		return
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn.Exec(
		ctx,
		"DELETE FROM lots WHERE batch_id = $1",
		id,
	)
}
