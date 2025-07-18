package dao

type SlugIterationTarget string

func (entity SlugIterationTarget) String() string {
	return string(entity)
}

const (
	SlugIterationTargetStoryPlan SlugIterationTarget = "story_plan"
	SlugIterationTargetLogline   SlugIterationTarget = "logline"
)
