package workflow

import (
	"fmt"
	"strings"
)

// RequiredFieldError represents an error that occurs when a required field is missing.
type RequiredFieldError struct {
	Field string
}

func (rf RequiredFieldError) Error() string {
	return fmt.Sprintf("'%s' is required but not provided", rf.Field)
}

// ValidationError represents an error that occurs when a field's value does not meet the expected constraints.
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

// Add appends a new error to the [ValidationError].
func (v *ValidationError) Add(err error) {
	v.errs = append(v.errs, err)
}

// HasError checks if the [ValidationError] contains any errors.
// Returns true if there are errors, otherwise returns false.
func (v *ValidationError) HasError() bool {
	return len(v.errs) > 0
}
