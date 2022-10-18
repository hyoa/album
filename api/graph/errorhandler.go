package graph

import (
	"errors"

	"github.com/hyoa/album/api/internal/translator"
	"github.com/hyoa/album/api/internal/user"
)

func HandleError(err error, translator translator.Translator) error {
	var messageId string
	data := make(map[string]interface{})
	switch e := err.(type) {
	case user.ErrorWithTranslation:
		messageId = e.ID()
	default:
		messageId = "InternalError"
	}

	return errors.New(translator.Translate(messageId, data))
}
