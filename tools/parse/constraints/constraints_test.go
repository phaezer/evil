package constraints

import (
	"github.com/stretchr/testify/assert"
	"go/build/constraint"
	"slices"
	"testing"
)

func TestIterBuildTagsInFile(t *testing.T) {
	src := `//go:build windows && !darwin
//go:build !linux && !windows
//go:build linux

package main

func main() {
	// comment that shouldn't be matched
}
`

	tags := slices.Collect(IterBuildTagsInFile("test.go", src, nil))
	assert.Equal(t, 3, len(tags))

	got := tags[0].String()
	want := "windows && !darwin"
	assert.Equal(t, want, got)

	got = tags[1].String()
	want = "!linux && !windows"
	assert.Equal(t, want, got)

	assert.True(t, tags[0].Eval(func(tag string) bool {
		return tag == "windows"
	}))

	assert.True(t, tags[1].Eval(func(tag string) bool {
		return tag == "darwin"
	}))

	assert.True(t, tags[2].Eval(func(tag string) bool {
		return tag == "linux"
	}))

	filteredTags := IterBuildTagsInFile("test.go", src, func(tag constraint.Expr) bool {
		return tag.Eval(func(tag string) bool {
			return tag == "linux"
		})
	})

	assert.Equal(t, 1, len(slices.Collect(filteredTags)))
}
