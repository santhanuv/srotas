package store

import (
	"maps"
)

// Store provides a dynamic storage mechanism for variables during config execution.
// It allows setting and retrieving variables as needed throughout execution.
type Store struct {
	variables map[string]any // variables contained in the Store.
}

// NewStore initializes and returns a new [Store] with the given variables.
func NewStore(ivars map[string]any) *Store {
	variables := make(map[string]any)

	maps.Copy(variables, ivars)

	return &Store{
		variables: variables,
	}
}

// Add merges the given variables into the Store.
// If a variable with the same name already exists, it will be overwritten.
func (s *Store) Add(vars map[string]any) {
	maps.Copy(s.variables, vars)
}

// Set sets the variable in the Store with given name and value.
// If a variable with the same name already exists, it will be overwritten.
func (s *Store) Set(name string, value any) {
	s.variables[name] = value
}

// Get retrieves the value of the variable identified by the given name,
// along with a boolean indicating whether the variable exists.
func (s *Store) Get(name string) (any, bool) {
	val, ok := s.variables[name]
	return val, ok
}

// Remove removes the variable identified by the given name.
func (s *Store) Remove(name string) {
	delete(s.variables, name)
}

// Map returns a map of all variables, where the keys are variable names and the values are their respective values.
func (s *Store) Map() map[string]any {
	return s.variables
}
