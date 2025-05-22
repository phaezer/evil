package evil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetaGenerator(t *testing.T) {
	mg := NewMetaGenerator("examples/simple/fib.go")

	assert.Equal(t, 1, len(mg.funcs))
}
