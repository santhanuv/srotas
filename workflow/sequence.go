package workflow

type Sequence struct {
	Name        string
	Description string
	Variables   map[string]any
	Steps       StepList
}
