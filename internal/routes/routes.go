package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", RenderIndex)

	fs := http.FileServer(http.Dir("web/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	RegisterContactRoutes(router)
	RegisterQuoteRoutes(router)

	// 404 & 5xx pages
	router.HandleFunc("GET /", RenderNotFound)
	router.HandleFunc("GET /505", RenderError)

	return router
}

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	templ, err := template.New("layout.html").ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
		"web/templates/sections/gallery.html",
		"web/templates/sections/amenities.html",
		"web/templates/sections/plan.html",
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
		return
	}

	err = templ.Execute(w, nil)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, 500, ErrorParams{})
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
