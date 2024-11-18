package config

type Definition struct {
	Version   string
	BaseUrl   string `yaml:"base_url"`
	Timeout   int
	MaxRetry  int               `yaml:"max_retry"`
	Variables map[string]any    `yaml:",flow"`
	Headers   map[string]string `yaml:",flow"`
	Sequence  Sequence
}
