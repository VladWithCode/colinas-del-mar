package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/vladwithcode/colinas/internal/db"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", RenderIndex)
	router.HandleFunc("GET /campana-fb", RenderCampana)

	fs := http.FileServer(http.Dir("web/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	RegisterContactRoutes(router)
	RegisterLotRoutes(router)
	RegisterQuoteRoutes(router)

	// 404 & 5xx pages
	router.HandleFunc("GET /", RenderNotFound)
	router.HandleFunc("GET /505", RenderError)

	return router
}

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	lots, err := db.GetLots()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
		return
	}

	templ, err := template.New("layout.html").ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
		"web/templates/contact-form.html",
		"web/templates/sections/gallery.html",
		"web/templates/sections/amenities.html",
		"web/templates/sections/plan.html",
		"web/templates/sections/green-area.html",
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
		return
	}

	lotMap := map[string]db.Lot{}

	for _, lot := range lots {
		id := fmt.Sprintf("Lt-%v_Mz-%v", lot.Lt, lot.Mz)

		lotMap[id] = *lot
	}

	err = templ.Execute(w, map[string]any{
		"Lot": map[string]any{
			"PriceCash":   172800,
			"PriceCredit": 192000,
			"Lt":          1,
			"Mz":          1,
			"Available":   true,
		},
		"Lots": lotMap,
	})

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
	}
}

func RenderCampana(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles(
		"web/templates/campana/campana.html",
		"web/templates/campana/campana_header.html",
		"web/templates/campana/campana_contact.html",
		"web/templates/campana/campana_feature_slider.html",
		"web/templates/campana/campana_footer.html",
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Ocurrio un error inesperado en el servidor."))
		return
	}

	err = templ.Execute(w, nil)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Ocurrio un error inesperado en el servidor."))
		return
	}
}

func RenderNotFound(w http.ResponseWriter, r *http.Request) {
}

func RenderError(w http.ResponseWriter, r *http.Request) {
}

func RespondWithError(w http.ResponseWriter, code int, params ErrorParams) {
	templ, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/error.html",
	)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ocurrio un error inesperado en el servidor."))
		return
	}

	if params.Message == "" {
		params.Message = "Ocurrió un error inesperado en el servidor."
	}
	if params.Code == 0 {
		params.Code = code
	}

	err = templ.Execute(w, params)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ocurrio un error inesperado en el servidor."))
		return
	}
}

type ErrorParams struct {
	Message string
	Code    int
}
