package alexa

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed data/active.*.toml
var localeFS embed.FS

type i18nLocalizer interface {
	localize(params localizeParams) string
}

type localizeParams struct {
	key         string
	data        map[string]interface{}
	pluralCount interface{}
}

func newLocalizer(userLocale string) i18nLocalizer {
	goi18nLocalizer := new(goi18nLocalizer)
	bundle := i18n.NewBundle(language.Spanish)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, _ = bundle.LoadMessageFileFS(localeFS, "data/active.es.toml")
	localizer := i18n.NewLocalizer(bundle, userLocale)
	goi18nLocalizer.localizer = localizer
	return goi18nLocalizer
}

type goi18nLocalizer struct {
	localizer *i18n.Localizer
}

func (l goi18nLocalizer) localize(params localizeParams) string {
	message, _ := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    params.key,
		TemplateData: params.data,
		PluralCount:  params.pluralCount,
	})
	return message
}
