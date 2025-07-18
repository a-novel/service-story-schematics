package models

type Slug string

func (slug Slug) String() string {
	return string(slug)
}
