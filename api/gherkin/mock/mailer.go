package mock

import (
	"encoding/json"
	"io/ioutil"
)

type mail struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Mailer struct{}

func (m *Mailer) SendMail(email, subject, body string) error {
	f, _ := json.Marshal(mail{Email: email, Subject: subject, Body: body})

	return ioutil.WriteFile("/tmp/mail.json", f, 0644)
}
