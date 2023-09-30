package gherkin_context

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/hyoa/album/api/internal/user"
)

type testMailerKey struct{}

type mail struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func ICheckInTheMailbox(ctx context.Context) (context.Context, error) {
	f, _ := ioutil.ReadFile("/tmp/mail.json")
	var mail mail

	json.Unmarshal(f, &mail)

	return context.WithValue(ctx, testMailerKey{}, mail), nil
}

func IShouldHaveAMailThatContainAPasswordResetLinkForWithSubject(ctx context.Context, email, subject string) error {
	res, _ := ctx.Value(testMailerKey{}).(mail)

	var re = regexp.MustCompile(`[\w]+:[\d]+\/resetPassword\?token=(?P<token>[\w=_\.-]+)`)
	matches := re.FindStringSubmatch(res.Body)

	if len(matches) != 2 {
		return errors.New("body does not contain any token")
	}

	tokenizer := user.CreateResetTokenizer()
	token, errDecode := tokenizer.Decode(matches[1])

	if res.Email != email || res.Subject != subject || errDecode != nil || token.Email != email {
		return errors.New("mail is not valid")
	}

	return nil
}

func IShouldHaveAMailThatContainAnInvitationLinkForWithSubject(ctx context.Context, email, subject string) error {
	res, _ := ctx.Value(testMailerKey{}).(mail)

	if res.Email != email || res.Subject != subject || res.Body != "Bonjour, nous vous invitons à voir nos albums photo à l'adresse suivante : localhost:3118" {
		return errors.New("mail is not valid")
	}

	return nil
}
