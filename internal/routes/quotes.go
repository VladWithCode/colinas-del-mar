package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vladwithcode/colinas/internal"
	"github.com/vladwithcode/colinas/internal/db"
	"github.com/vladwithcode/colinas/internal/whatsapp"
)

func RegisterQuoteRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /api/quote", CreateNewQuote)
}

func CreateNewQuote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/layout.html")

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.WriteHeader(500)
		err = templ.ExecuteTemplate(
			w,
			"quote-form",
			map[string]string{
				"ErrorMessage": "El formulario contiene información inválida",
			},
		)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}
	}

	var (
		name  = r.Form.Get("name")
		phone = r.Form.Get("phone")
	)

	date := time.Now()
	phone, err = internal.FormatPhone(phone)

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, _ := template.ParseFiles("web/templates/layout.html")

		w.WriteHeader(400)
		err = templ.ExecuteTemplate(
			w,
			"contact-form",
			map[string]string{
				"ErrorMessage": "El número de teléfono es inválido",
			},
		)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}
	}

	id, _ := uuid.NewV7()

	contact := db.Contact{
		Id:        id.String(),
		Name:      name,
		Phone:     phone,
		CreatedAt: date,
		Pending:   true,
	}

	err = whatsapp.SendNewContactNotification(&contact)

	if err != nil {
		templ, _ := template.ParseFiles("web/templates/layout.html")

		w.WriteHeader(500)
		err = templ.ExecuteTemplate(
			w,
			"contact-form",
			map[string]string{
				"ErrorMessage": "Ocurrió un error al enviar tu solicitud. Intenta de nuevo más tarde",
			},
		)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}
	}

	w.WriteHeader(201)
	templ, err := template.ParseFiles("web/templates/layout.html")

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}

	err = templ.ExecuteTemplate(
		w,
		"contact-form",
		map[string]string{
			"SuccessMessage": "Se envió tu solicitud. En breve nos comunicaremos contigo.",
		},
	)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}
}
