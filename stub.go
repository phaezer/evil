package evil

import "fmt"

// Stub is used as a reference during code generation
func Stub[T any](v string) T {
	// stub only used as a reference during code generation
	panic(fmt.Errorf("stub was called: %s", v))
}
