package utils

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func SetFromSlice[T comparable](slice []T) *Set[T] {
	set := &Set[T]{items: make(map[T]struct{})}
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

func (s *Set[T]) Add(value T) *Set[T] {
	s.items[value] = struct{}{}
	return s
}

func (s *Set[T]) Remove(value T) *Set[T] {
	delete(s.items, value)
	return s
}

func (s *Set[T]) Has(value T) bool {
	_, exists := s.items[value]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Values() []T {
	values := make([]T, 0)
	for key := range s.items {
		values = append(values, key)
	}
	return values
}

func (s *Set[T]) Equal(set *Set[T]) bool {
	if len(s.items) != len(set.items) {
		return false
	}

	for key := range set.items {
		if !s.Has(key) {
			return false
		}
	}
	return true
}
