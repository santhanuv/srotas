package store

type Store struct {
	variables map[string]any
}

func NewStore(initialValues map[string]any) *Store {
	variables := make(map[string]any)

	if initialValues != nil {
		for key, val := range initialValues {
			variables[key] = val
		}
	}

	return &Store{
		variables: variables,
	}
}

func (s *Store) Add(vars map[string]any) {
	for vn, vv := range vars {
		s.variables[vn] = vv
	}
}

func (s *Store) Set(key string, value any) {
	s.variables[key] = value
}

func (s *Store) Get(key string) (any, bool) {
	val, ok := s.variables[key]
	return val, ok
}

func (s *Store) Remove(key string) {
	delete(s.variables, key)
}

func (s *Store) ToMap() map[string]any {
	return s.variables
}
