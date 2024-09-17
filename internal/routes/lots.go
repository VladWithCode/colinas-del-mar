package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vladwithcode/colinas/internal"
	"github.com/vladwithcode/colinas/internal/db"
)

func RegisterLotRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /api/lots", RegisterLots)
	r.HandleFunc("GET /api/lots/{lt}/{mz}", GetLotByLtMz)
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

func GetLotByLtMz(w http.ResponseWriter, r *http.Request) {
	lt := r.PathValue("lt")
	mz := r.PathValue("mz")

	lot, err := db.GetLotByLtMz(lt, mz)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{Message: "Lot Id not found"})
		return
	}

	templ, err := template.ParseFiles("web/templates/sections/plan.html")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{Message: "Error parsing template"})
		return
	}

	err = templ.ExecuteTemplate(
		w,
		"popup-content",
		map[string]any{
			"Lot":           lot,
			"PriceCash":     internal.FormatMoney(lot.PriceCash),
			"PriceMtCash":   internal.FormatMoney(lot.PriceCash / lot.Area),
			"PriceCredit":   internal.FormatMoney(lot.PriceCredit),
			"PriceMtCredit": internal.FormatMoney(lot.PriceCredit / lot.Area),
		},
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{Message: "Error executing template"})
		return
	}
}
