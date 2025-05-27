package evil

import (
	"go/ast"
	"iter"
)

// IterComments returns a comment iterator for all comment lines in a [ast.File] node
//
// if filter is nil, all comments are returned
// if filter is not nil, only comments that pass the filter are returned
func IterComments(node *ast.File, f func(line string) bool) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, commentGroup := range node.Comments {
			for _, comment := range commentGroup.List {
				text := comment.Text
				if f == nil || f(text) {
					if !yield(text) {
						return
					}
				}
			}
		}
	}
}
