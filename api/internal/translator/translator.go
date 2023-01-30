package translator

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Translator struct {
	localizer *i18n.Localizer
}

func CreateTranslator(path string) *Translator {
	bundle := i18n.NewBundle(language.French)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile(path)

	return &Translator{
		localizer: i18n.NewLocalizer(bundle, "fr"),
	}
}

func (t *Translator) Translate(messageId string, data map[string]interface{}) string {
	message, _ := t.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: data,
	})

	return message
}
