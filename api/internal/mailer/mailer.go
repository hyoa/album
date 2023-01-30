package mailer

import (
	"os"

	"github.com/hyoa/album/api/internal/translator"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridMailer struct {
	ApiKey     string
	Translater translator.Translator
}

func (sgm *SendgridMailer) SendMail(email, subjectKey, bodyKey string, bodyData map[string]interface{}) error {
	from := mail.NewEmail(os.Getenv("APP_NAME"), os.Getenv("ADMIN_EMAIL"))

	to := mail.NewEmail(email, email)

	subjectTranslated := sgm.Translater.Translate(subjectKey, bodyData)
	bodyTranslated := sgm.Translater.Translate(bodyKey, bodyData)

	message := mail.NewSingleEmail(from, subjectTranslated, to, bodyTranslated, "")
	client := sendgrid.NewSendClient(sgm.ApiKey)
	_, err := client.Send(message)

	return err
}
