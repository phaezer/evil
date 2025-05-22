package comments

import (
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"slices"
)

// IterInFile returns a comment iterator for all comment lines in a file
//
// filePath and src are passed to [parser.ParseFile] for parsing
// if [parser.ParseFiles] returns an error, IterInFile panics
//
// if filter is nil, all comments are returned
// if filter is not nil, only comments that pass the filter are returned
func IterInFile(filePath string, src any, filter func(line string) bool) iter.Seq[string] {
	if filter == nil {
		filter = func(_ string) bool { return true }
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return IterNode(node, filter)
}

// IterNode returns a comment iterator for all comment lines in a [ast.File] node
//
// if filter is nil, all comments are returned
// if filter is not nil, only comments that pass the filter are returned
func IterNode(node *ast.File, filter func(line string) bool) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, commentGroup := range node.Comments {
			for _, comment := range commentGroup.List {
				text := comment.Text
				if !filter(text) {
					continue
				}

				if !yield(text) {
					return
				}
			}
		}
	}
}

// Collect returns a slice of all comment lines in a file
//
// filePath and src are passed to [parser.ParseFile] for parsing
// if [parser.ParseFiles] returns an error, IterInFile panics
//
// if filter is nil, all comments are returned
// if filter is not nil, only comments that pass the filter are returned
func Collect(filePath string, src any, filter func(line string) bool) []string {
	return slices.Collect(IterInFile(filePath, src, filter))
}
