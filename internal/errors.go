package internal

import (
	"fmt"
	"strings"
)

type RequiredFieldError struct {
	Field string
}

func (rf RequiredFieldError) Error() string {
	return fmt.Sprintf("'%s' is required but not provided", rf.Field)
}

type ValidationError struct {
	errs []error
}

func (v *ValidationError) Error() string {
	var messages []string

	for _, err := range v.errs {
		messages = append(messages, err.Error())
	}

	joinedMessage := strings.Join(messages, "\n\t")
	return fmt.Sprintf("validation errors:\n\t%s;", joinedMessage)
}

func (v *ValidationError) Add(err error) {
	v.errs = append(v.errs, err)
}

func (v *ValidationError) HasError() bool {
	return len(v.errs) > 0
}
