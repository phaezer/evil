package simple

//go:generate evil

// Fibonacci function returns a slice containing the Fibonacci sequence up to n terms.
func Fibonacci(n int) []int {
	if n <= 0 {
		return []int{}
	}
	if n == 1 {
		return []int{0}
	}

	fib := make([]int, n)
	fib[0] = 0
	fib[1] = 1

	for i := 2; i < n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	return fib
}

var fib []int
