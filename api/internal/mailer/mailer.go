package mailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridMailer struct {
	ApiKey string
}

func (sgm *SendgridMailer) SendMail(email, subject, body string) error {
	from := mail.NewEmail("Pauline&Jules", "admin@pauline-jules.fr")

	to := mail.NewEmail(email, email)

	fmt.Println(from, to, body, subject)
	message := mail.NewSingleEmail(from, subject, to, body, "")
	client := sendgrid.NewSendClient(sgm.ApiKey)
	r, err := client.Send(message)
	fmt.Println(r)

	return err
}
