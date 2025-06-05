package counter

import (
	"context"
	"log"
	"os"

	"github.com/vladwithcode/colinas/internal/db"
)

type VisitCounter struct {
	Id     string `json:"id" db:"id"`
	Route  string `json:"route" db:"route"`
	Visits int    `json:"visits" db:"visits"`
}

func InsertOrUpdateRouteVisit(forRoute string) error {
	conn, err := db.GetPool()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		`INSERT INTO visit_counter (id, route, visits) 
			VALUES (gen_random_uuid(), $1, 1) 
			ON CONFLICT (route) DO UPDATE SET visits = visit_counter.visits + 1`,
		forRoute,
	)

	if err != nil {
		return err
	}

	return nil
}

func RegisterVisit(route string) {
	go func() {
		err := InsertOrUpdateRouteVisit(route)
		if err != nil {
			fpath := os.ExpandEnv("${HOME}/.local/log/colinas/error.log")
			visitLogFile, oErr := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
			if oErr != nil {
				log.Printf("Error opening visit log: %v\n", oErr)
				return
			}
			defer visitLogFile.Close()

			logger := log.New(visitLogFile, "VisitLog error: ", log.LstdFlags)
			logger.Printf("for route %s\n\t%v", route, err)
		}
	}()
}
