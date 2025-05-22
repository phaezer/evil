//go:build evil

package main

import (
	"math/big"

	"github.com/phaezer/evil"
	"github.com/phaezer/evil/examples/simple"
)

// pi code credit to https://github.com/JJ/pigo and https://go.dev/play/p/mb5eoZpnYN
func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x * x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false

	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}

func Pi(ndigits int64) string {
	if ndigits <= 7 {
		return "3.141595"
	} else {
		digits := big.NewInt(ndigits + 10)
		unity := big.NewInt(0)                 // crea un entero tocho
		unity.Exp(big.NewInt(10), digits, nil) // Le asigna valor
		pi := big.NewInt(0)
		four := big.NewInt(4) // Todos deben ser enteros tocho

		// Serie de McLaurin
		pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
		output := fmt.Sprintf("%s.%s", pi.String()[0:1], pi.String()[1:ndigits])
		return output
	}
}

func main() {
	g := evil.NewGenerator()

	defer g.generate()

	piNumDigits := 100
	fib := simple.Fibonacci(1000)
	pi100 := Pi(piNumDigits)

	// initialize a variable in the global scope using an function from the package
	g.InitVar(g, "fib", fib)
	// create a const named pi
	g.NewConst("pi", pi100)
	// create a var named PiVar
	g.NewVar("PiVar", pi100)

	// create a function named getPi
	g.NewFunc("getPi", func(n int) string {
		if n > len(pi100) {
			return pi100
		}

		// use evil.Stub that will be swapped with the string contents during code generation
		return evil.Stub[string]("pi")[:n]
	})

	// create a function with string
	g.Raw(`
func getFib(n int)) []int {
	if n > len(fib) {
		return fib
	}
	return fib[n:]
}
`)
}
