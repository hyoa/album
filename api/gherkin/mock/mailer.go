package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hyoa/album/api/internal/translator"
)

type mail struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Mailer struct {
	translator.Translator
}

func (m *Mailer) SendMail(email, subjectKey, bodyKey string, bodyData map[string]interface{}) error {
	subjectTranslated := m.Translator.Translate(subjectKey, bodyData)
	bodyTranslated := m.Translator.Translate(bodyKey, bodyData)

	fmt.Println(subjectTranslated, bodyTranslated)
	f, _ := json.Marshal(mail{Email: email, Subject: subjectTranslated, Body: bodyTranslated})

	return ioutil.WriteFile("/tmp/mail.json", f, 0644)
}
