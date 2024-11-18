package validation

type Validation interface {
	Validate() error
}

type StatusValidation struct {
	Type   string
	Status int
}
