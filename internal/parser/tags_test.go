package parser

import (
	"github.com/stretchr/testify/assert"
	"go/build/constraint"
	"go/token"
	"slices"
	"testing"
)

const testSrc = `//go:build windows && !darwin
//go:build !linux && !windows
//go:build linux

package main

func main() {
	// comment in main
}
`

func TestIterBuildTagsInFile(t *testing.T) {
	fs := token.NewFileSet()

	tags := slices.Collect(IterInFile(fs, "test.go", testSrc, nil))
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

	filteredTags := IterInFile(fs, "test.go", testSrc, func(tag constraint.Expr) bool {
		return tag.Eval(func(tag string) bool {
			return tag == "linux"
		})
	})

	assert.Equal(t, 1, len(slices.Collect(filteredTags)))
}

func TestIterWithAnyBuildTags(t *testing.T) {
	fs := token.NewFileSet()
	tags := slices.Collect(IterInFile(fs, "test.go", testSrc, WithAny("linux", "windows")))
	assert.Equal(t, 2, len(tags))

}

func TestIterWithAllBuildTags(t *testing.T) {
	fs := token.NewFileSet()
	tags := slices.Collect(IterInFile(fs, "test.go", testSrc, WithAll("darwin", "arm64")))
	assert.Equal(t, 1, len(tags))

}
