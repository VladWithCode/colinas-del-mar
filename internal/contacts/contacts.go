package contacts

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

func LogContactRequest(name, phone string, date time.Time) {
	logPath, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	logPath = logPath + "/.local/log/colinas_sibra/contact.log"

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer logFile.Close()

	templ, err := template.ParseFiles("internal/contacts/conreq_templ")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	templ.Execute(logFile, map[string]string{
		"Name":  name,
		"Phone": phone,
		"Date":  date.Format("2006-01-02"),
	})
}
