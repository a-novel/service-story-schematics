package models

type Lang string

func (lang Lang) String() string {
	return string(lang)
}

const (
	LangEN Lang = "en"
	LangFR Lang = "fr"
)
