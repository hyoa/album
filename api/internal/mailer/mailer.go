package mailer

import (
	"github.com/hyoa/album/api/internal/translator"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridMailer struct {
	ApiKey     string
	Translater translator.Translator
}

func (sgm *SendgridMailer) SendMail(email, subjectKey, bodyKey string, bodyData map[string]interface{}) error {
	from := mail.NewEmail("Pauline&Jules", "admin@pauline-jules.fr")

	to := mail.NewEmail(email, email)

	subjectTranslated := sgm.Translater.Translate(subjectKey, bodyData)
	bodyTranslated := sgm.Translater.Translate(bodyKey, bodyData)

	message := mail.NewSingleEmail(from, subjectTranslated, to, bodyTranslated, "")
	client := sendgrid.NewSendClient(sgm.ApiKey)
	_, err := client.Send(message)

	return err
}
