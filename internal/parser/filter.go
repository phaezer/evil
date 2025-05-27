package parser

func And[V any](fns ...func(V) bool) func(V) bool {
	return func(item V) bool {
		for _, f := range fns {
			if !f(item) {
				return false
			}
		}

		return true
	}
}

func Or[V any](fns ...func(V) bool) func(V) bool {
	return func(item V) bool {
		for _, f := range fns {
			if f(item) {
				return true
			}
		}

		return false
	}
}
