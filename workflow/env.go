package workflow

type Env struct {
	Variables map[string]any
	Headers   map[string][]string
}

func (e *Env) AppendVars(varsList ...map[string]any) {
	for _, v := range varsList {
		for name, val := range v {
			e.Variables[name] = val
		}
	}
}

func (e *Env) AppendHeaders(headersList ...map[string][]string) {
	for _, headers := range headersList {
		for key, val := range headers {
			e.Headers[key] = append(e.Headers[key], val...)
		}
	}
}

func NewEnv(variables map[string]any, headers map[string][]string) *Env {
	if variables == nil {
		variables = make(map[string]any)
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	return &Env{
		Variables: variables,
		Headers:   headers,
	}
}
