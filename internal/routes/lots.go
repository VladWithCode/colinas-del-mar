package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vladwithcode/colinas/internal/db"
)

func RegisterLotRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /api/lots", RegisterLots)
}

func RegisterLots(w http.ResponseWriter, r *http.Request) {
	var lotsData []map[string]any

	err := json.NewDecoder(r.Body).Decode(&lotsData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
		return
	}

	lotBatch := pgx.Batch{}
	batchId, _ := uuid.NewV7()
	for _, lot := range lotsData {
		newLot := db.Lot{}
		id, _ := uuid.NewV7()
		newLot.Id = id.String()
		newLot.Lt = lot["lt"].(string)
		newLot.Mz = lot["mz"].(string)
		newLot.Area = lot["area"].(float64)
		newLot.Measures = lot["measures"].(string)
		newLot.Type = lot["type"].(string)
		newLot.PriceCash = lot["priceCash"].(float64)
		newLot.PriceCredit = lot["priceCredit"].(float64)
		newLot.BatchId = batchId.String()

		db.QueueCreateLot(&lotBatch, &newLot)
	}

	err = db.BatchCreateLot(&lotBatch)

	if err != nil {
		db.RemoveFailedBatchLots(batchId.String())
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("OK"))
}
