package seq

import "iter"

func Once[T any](t T) iter.Seq[T] {
	return func(yield func(T) bool) {
		if !yield(t) {
			return
		}
	}
}

func Once2[T1, T2 any](t1 T1, t2 T2) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		if !yield(t1, t2) {
			return
		}
	}
}
