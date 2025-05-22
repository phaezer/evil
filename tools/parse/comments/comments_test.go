package comments

import (
	"github.com/stretchr/testify/assert"
	"log"
	"slices"
	"testing"
)

func TestIterCommentsInFile(t *testing.T) {
	src := `
package main

// comment 1
// comment 2
// comment 3

func main() {
	// comment 4
	// comment 5
	// comment 6
}
`

	comments := slices.Collect(IterInFile("test.go", src, nil))
	for _, comment := range comments {
		log.Printf("comment = %#v\n", comment)
	}
	assert.Len(t, comments, 6)
}
