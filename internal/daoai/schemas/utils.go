package schemas

type Schema struct {
	Description string `yaml:"description"`
	Schema      any    `yaml:"schema"`
}
