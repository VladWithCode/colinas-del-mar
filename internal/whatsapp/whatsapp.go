package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/vladwithcode/colinas/internal/db"
)

type TemplateVar map[string]any

type TemplateComponent struct {
	ComponentType string        `json:"type"`
	Parameters    []TemplateVar `json:"parameters"`
}

type TemplateData struct {
	TemplateName string
	BodyVars     []TemplateVar
	HeaderVars   []TemplateVar
}

type templatePayload struct {
	Name       string              `json:"name"`
	Components []TemplateComponent `json:"components"`
	Language   struct {
		Code string `json:"code"`
	} `json:"language"`
}

var (
	ErrPhoneNotSet = errors.New("phone number is not set")
)

func postCloudAPIMessage(requestPayload any) error {
	phoneNumberId := os.Getenv("FB_PHONE_NUMBER_ID")
	fbAccessToken := os.Getenv("FB_ACCESS_TOKEN")

	if phoneNumberId == "" {
		panic("FB_PHONE_NUMBER_ID not in ENV")
	}
	if fbAccessToken == "" {
		panic("FB_ACCESS_TOKEN not in ENV")
	}

	reqBody, err := json.Marshal(requestPayload)

	if err != nil {
		return err
	}

	reqUrl := fmt.Sprintf("https://graph.facebook.com/v22.0/%v/messages", phoneNumberId)
	req, err := http.NewRequest("post", reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", fbAccessToken))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("Errored with status: %v\n", resp.Status)
		fmt.Printf("resp.Header: %v\n", resp.Header["Www-Authenticate"])

		return fmt.Errorf("request failed with code: %v", resp.StatusCode)
	}

	return nil
}

func SendTemplateMessage(phoneNumber string, data TemplateData) error {
	var reqPayload struct {
		MessagingProduct string          `json:"messaging_product"`
		MessageType      string          `json:"type"`
		ToPhone          string          `json:"to"`
		Template         templatePayload `json:"template"`
	}

	var components []TemplateComponent

	if len(data.HeaderVars) != 0 {
		components = append(components, TemplateComponent{
			ComponentType: "header",
			Parameters:    data.HeaderVars,
		})
	}

	if len(data.BodyVars) != 0 {
		components = append(components, TemplateComponent{
			ComponentType: "body",
			Parameters:    data.BodyVars,
		})
	}

	reqPayload.MessageType = "template"
	reqPayload.MessagingProduct = "whatsapp"
	reqPayload.ToPhone = phoneNumber
	reqPayload.Template = templatePayload{
		Name: data.TemplateName,
		Language: struct {
			Code string `json:"code"`
		}{Code: "es"},
		Components: components,
	}

	return postCloudAPIMessage(reqPayload)
}

func SendNewContactNotification(contactData *db.Contact) error {
	sendToPhone := os.Getenv("CONTACT_PHONE")

	if sendToPhone == "" {
		return ErrPhoneNotSet
	}

	var (
		headerVars []TemplateVar
		bodyVars   []TemplateVar
	)

	createdDateStr := contactData.CreatedAt.Format("2006-01-02 03:04pm")
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": contactData.Name,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": createdDateStr,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type": "text",
		"text": contactData.Phone,
	})

	return SendTemplateMessage(sendToPhone, TemplateData{
		TemplateName: "info_request",
		BodyVars:     bodyVars,
		HeaderVars:   headerVars,
	})
}

func SendNewCampaignNotification(contactData *db.Contact) error {
	sendToPhone := os.Getenv("CAMPAIGN_PHONE")

	if sendToPhone == "" {
		return ErrPhoneNotSet
	}

	var (
		headerVars []TemplateVar
		bodyVars   []TemplateVar
	)

	createdDateStr := contactData.CreatedAt.Format("2006-01-02 03:04pm")
	bodyVars = append(bodyVars, TemplateVar{
		"type":           "text",
		"parameter_name": "name",
		"text":           contactData.Name,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type":           "text",
		"parameter_name": "date",
		"text":           createdDateStr,
	})
	bodyVars = append(bodyVars, TemplateVar{
		"type":           "text",
		"parameter_name": "phone",
		"text":           contactData.Phone,
	})

	return SendTemplateMessage(sendToPhone, TemplateData{
		TemplateName: "campaing_request",
		BodyVars:     bodyVars,
		HeaderVars:   headerVars,
	})
}
