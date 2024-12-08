package config

import "github.com/santhanuv/srotas/config/step/http"

type Definition struct {
	Version   string
	BaseUrl   string `yaml:"base_url"`
	Timeout   uint
	MaxRetry  uint           `yaml:"max_retry"`
	Variables map[string]any `yaml:",flow"`
	Headers   http.Header    `yaml:",flow"`
	Sequence  Sequence
}
