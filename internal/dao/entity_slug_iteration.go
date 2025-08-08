package dao

type SlugIterationTarget string

func (entity SlugIterationTarget) String() string {
	return string(entity)
}

const (
	SlugIterationTargetLogline SlugIterationTarget = "logline"
)
