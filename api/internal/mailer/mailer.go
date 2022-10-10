package mailer

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridMailer struct {
	ApiKey string
}

func (sgm *SendgridMailer) SendMail(email, subject, body string) error {
	from := mail.NewEmail("Pauline&Jules", "admin@pauline-jules.fr")

	to := mail.NewEmail(email, email)

	message := mail.NewSingleEmail(from, subject, to, body, "")
	client := sendgrid.NewSendClient(sgm.ApiKey)
	_, err := client.Send(message)

	return err
}
