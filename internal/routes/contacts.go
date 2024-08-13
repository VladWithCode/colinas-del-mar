package routes

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/vladwithcode/colinas/internal/db"
	"github.com/vladwithcode/colinas/internal/whatsapp"
)

func RegisterContactRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/contact", CreateNewContact)
}

func CreateNewContact(w http.ResponseWriter, r *http.Request) {
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
			"contact-form",
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
	phone, err = formatPhone(phone)

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

func formatPhone(p string) (string, error) {
	stripCountryCodeExp := regexp.MustCompile(`^\+52`)
	replaceExp := regexp.MustCompile(`[ -]`)
	numExp := regexp.MustCompile(`[0-9]{10}`)

	phone := stripCountryCodeExp.ReplaceAll([]byte(p), []byte(""))
	phone = replaceExp.ReplaceAll([]byte(phone), []byte(""))

	if !numExp.Match(phone) {
		return "", errors.New(fmt.Sprintf("The string is not a valid phone number: %v", p))
	}

	return string(phone), nil
}
