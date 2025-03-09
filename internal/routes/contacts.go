package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vladwithcode/colinas/internal"
	"github.com/vladwithcode/colinas/internal/contacts"
	"github.com/vladwithcode/colinas/internal/db"
	"github.com/vladwithcode/colinas/internal/whatsapp"
)

func RegisterContactRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /api/contact", CreateNewContact)
	router.HandleFunc("POST /api/contact/lot", CreateNewContact)
	router.HandleFunc("POST /api/contact/campana", CreateNewCampaingContact)
}

func CreateNewContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/contact-form.html")

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Add("X-ForceSwap", "true")
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
		}
		return
	}

	var (
		name  = r.Form.Get("name")
		phone = r.Form.Get("phone")
		date  = time.Now()
	)

	contacts.LogContactRequest(name, phone, date)

	phone, err = internal.FormatPhone(phone)

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, _ := template.ParseFiles("web/templates/contact-form.html")

		w.Header().Add("X-ForceSwap", "true")
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
		}
		return
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
		fmt.Printf("err: %v\n", err)
		templ, _ := template.ParseFiles("web/templates/contact-form.html")

		w.Header().Add("X-ForceSwap", "true")
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
		}
		return
	}

	templ, err := template.ParseFiles("web/templates/contact-form.html")

	if err != nil {
		fmt.Printf("err: %v\n", err)
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
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}
}

func CreateNewCampaingContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/campana/campana_contact.html")

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Add("X-ForceSwap", "true")
		w.Header().Add("HX-Reswap", "contact-feedback")
		w.WriteHeader(400)
		err = templ.ExecuteTemplate(
			w,
			"request-feedback",
			map[string]any{
				"IsInvalid": true,
				"Errors":    []string{"Ocurrió un error al procesar la información del formulario. Verifica tus datos e intenta de nuevo."},
			},
		)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}
		return
	}

	var (
		name  = r.Form.Get("name")
		phone = r.Form.Get("phone")
		date  = time.Now()
	)

	contacts.LogContactRequest(name, phone, date)

	phoneFmt, err := internal.FormatPhone(phone)

	if err != nil {
		fmt.Printf("Parse form err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/campana/campana_contact.html")
		if err != nil {
			fmt.Printf("Parse files err: %v\n", err)
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}

		w.Header().Add("X-ForceSwap", "true")
		w.WriteHeader(400)
		err = templ.ExecuteTemplate(
			w,
			"contact-form",
			map[string]any{
				"IsInvalid":      true,
				"PhoneInvalid":   true,
				"Errors":         []string{"El número de teléfono es inválido"},
				"PhoneValue":     phone,
				"NameValue":      name,
				"ConsentChecked": true,
			},
		)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}
		return
	}

	id := uuid.Must(uuid.NewV7())

	contact := db.Contact{
		Id:         id.String(),
		Name:       name,
		Phone:      phoneFmt,
		CreatedAt:  date,
		Pending:    true,
		CampaingId: "fb-mobile-campaign",
	}

	err = whatsapp.SendNewCampaignNotification(&contact)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		templ, err := template.ParseFiles("web/templates/campana/campana_contact.html")
		if err != nil {
			fmt.Printf("Parse form err: %v\n", err)
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
			return
		}

		w.Header().Add("X-ForceSwap", "true")
		w.Header().Add("Hx-Retarget", "#contact-feedback")
		w.WriteHeader(500)
		err = templ.ExecuteTemplate(
			w,
			"request-feedback",
			map[string]any{
				"IsInvalid":      true,
				"Errors":         []string{"Ocurrió un error al enviar tu solicitud. Intenta de nuevo más tarde"},
				"PhoneValue":     phone,
				"NameValue":      name,
				"ConsentChecked": true,
			},
		)

		if err != nil {
			fmt.Printf("Execute templ err: %v\n", err)
			RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		}
		return
	}

	err = db.CreateContact(&contact)
	if err != nil {
		log.Printf("error saving contact to db: %v\n", err)
	}

	templ, err := template.ParseFiles("web/templates/campana/campana_contact.html")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}

	err = templ.ExecuteTemplate(
		w,
		"contact-form",
		map[string]any{
			"IsSuccess":      true,
			"PhoneValue":     phone,
			"NameValue":      name,
			"ConsentChecked": true,
		},
	)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, ErrorParams{})
		return
	}
}

func CreateNewContactWithLot(w http.ResponseWriter, r *http.Request) {
}
