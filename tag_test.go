package evil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchTagsInFileNodeMultiLine(t *testing.T) {
	src := `
//go:build windows
//go:build linux && !darwin

package main

func main() {}
	`

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "test.go", strings.NewReader(src), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	match, err := MatchTagsInFileNode(file, []string{"linux", "windows"}, []string{"ignore"})
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, match)

	match, err = MatchTagsInFileNode(file, []string{"linux", "darwin"}, []string{"ignore"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)

	match, err = MatchTagsInFileNode(file, []string{"linux"}, []string{"ignore"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)

	match, err = MatchTagsInFileNode(file, []string{"linux"}, []string{"windows"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)

	match, err = MatchTagsInFileNode(file, []string{"linux", "windows"}, []string{"windows"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)
}

func TestMatchTagsInFileNodeSingleLine(t *testing.T) {
	src := `
//go:build (windows || linux) && !dev

package main

func main() {}
	`

	var (
		file  *ast.File
		match bool
		err   error
	)

	fs := token.NewFileSet()
	file, err = parser.ParseFile(fs, "test.go", strings.NewReader(src), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	match, err = MatchTagsInFileNode(file, []string{"windows", "prod"}, []string{"ignore"})
	if err != nil {
		t.Fatal(err)
	}
	match, err = MatchTagsInFileNode(file, []string{"linux", "darwin"}, []string{"ignore"})
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, match)

	match, err = MatchTagsInFileNode(file, []string{"linux"}, []string{"windows"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)

	match, err = MatchTagsInFileNode(file, nil, []string{"linux"})
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, match)
}
