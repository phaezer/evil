package filter

type Filter[T any] func(T) bool

func And[T any](filters ...Filter[T]) Filter[T] {
	return func(item T) bool {
		for _, filter := range filters {
			if !filter(item) {
				return false
			}
		}

		return true
	}
}

func Or[T any](filters ...Filter[T]) Filter[T] {
	return func(item T) bool {
		for _, filter := range filters {
			if filter(item) {
				return true
			}
		}

		return false
	}
}
