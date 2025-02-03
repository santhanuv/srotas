package workflow

type Definition struct {
	Version   string
	BaseUrl   string `yaml:"base_url"`
	Timeout   uint
	MaxRetry  uint `yaml:"max_retry"`
	Variables map[string]string
	Headers   Header
	Steps     StepList
	Output    map[string]string
	OutputAll bool `yaml:"output_all"`
}
