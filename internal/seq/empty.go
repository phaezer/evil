package seq

import "iter"

func Empty[K any]() iter.Seq[K] {
	return func(yield func(K) bool) {
		// do nothing as there are no values
	}
}

func Empty2[V, K any]() iter.Seq2[V, K] {
	return func(yield func(V, K) bool) {
		// do nothing as there are no values
	}
}
