package graph

import (
	"errors"

	"github.com/hyoa/album/api/internal/translator"
	"github.com/hyoa/album/api/internal/user"
)

func HandleError(err error, translator translator.Translator) error {
	var messageId string
	data := make(map[string]interface{})

	var (
		userError *user.UserError
	)

	if errors.As(err, &userError) {
		errI18N := err.(*user.UserError)
		messageId = errI18N.I18NID()
	} else {
		messageId = "InternalError"
	}

	return errors.New(translator.Translate(messageId, data))
}
