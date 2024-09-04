package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	Id         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Phone      string    `json:"phone" db:"phone"`
	Pending    bool      `json:"pending" db:"pending"`
	QuoteDate  time.Time `json:"quoteDate" db:"quote_date"`
	AttendedAt time.Time `json:"attendedDate" db:"attended_date"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	AssignedTo string    `json:"assignedTo" db:"assign_to"`
}

func CreateQuote(quote *Quote) error {
	conn, err := GetPool()
	if err != nil {
		return err
	}
	defer conn.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id, _ := uuid.NewV7()
	_, err = conn.Exec(
		ctx,
		"INSERT INTO quote (id, name, phone, pending, quote_date, attended_at, created_at) VALUES $1, $2, $3, $4, $5, $6, $7",
		id,
		quote.Name,
		quote.Phone,
		quote.Pending,
		quote.QuoteDate,
		quote.AttendedAt,
		quote.CreatedAt,
	)

	quote.Id = id.String()
	return nil
}
