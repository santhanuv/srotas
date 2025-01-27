package workflow

type Definition struct {
	Version   string
	BaseUrl   string `yaml:"base_url"`
	Timeout   uint
	MaxRetry  uint           `yaml:"max_retry"`
	Variables map[string]any `yaml:",flow"`
	Headers   Header         `yaml:",flow"`
	Sequence  Sequence
}
