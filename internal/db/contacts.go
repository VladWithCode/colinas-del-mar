package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Id         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Phone      string    `json:"phone" db:"phone"`
	Pending    bool      `json:"pending" db:"pending"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	AssignedTo string    `json:"assignedTo" db:"assign_to"`
}

func CreateContact(contact *Contact) error {
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
		"INSERT INTO contacts (id, name, phone, pending, created_at, assigned_to) VALUES $1, $2, $3, $4, $5, $6",
		id,
		contact.Name,
		contact.Phone,
		contact.Pending,
		contact.CreatedAt,
		contact.AssignedTo,
	)

	contact.Id = id.String()
	return nil
}
