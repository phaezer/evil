package constraints

import (
	"github.com/phaezer/evil/tools/filter"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"iter"
	"strings"

	"github.com/phaezer/evil/tools/parse/comments"
)

type ConstraintFilter filter.Filter[constraint.Expr]

func WithBuildTags(tags []string) ConstraintFilter {
	return func(c constraint.Expr) bool {
		var ok int
		for _, tagName := range tags {
			if c.Eval(func(t string) bool {
				return t == tagName
			}) {
				ok += 1
			}
		}

		return ok == len(tags)
	}
}

func WithAnyBuildTags(tags []string) ConstraintFilter {
	return func(c constraint.Expr) bool {
		for _, tagName := range tags {
			if c.Eval(func(t string) bool {
				return t == tagName
			}) {
				return true
			}
		}

		return false
	}
}

// IterBuildTagsInFile returns an iterator for all build tags in a file
//
// filePath and src are passed to [parser.ParseFile] for parsing
// if [parser.ParseFiles] returns an error, IterBuildTagsInFile panics
//
// if filter is nil, all build tags are returned
// if filter is not nil, only build tags that pass the filter are returned
func IterBuildTagsInFile(filePath string, src any, filter ConstraintFilter) iter.Seq[constraint.Expr] {
	if filter == nil {
		filter = func(expr constraint.Expr) bool { return true }
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return IterBuildTags(node, filter)
}

// IterBuildTags returns an iterator for all build tags in an [ast.File] node
//
// if filter is nil, all build tags are returned
// if filter is not nil, only build tags that pass the filter are returned
func IterBuildTags(node *ast.File, filter ConstraintFilter) iter.Seq[constraint.Expr] {
	return func(yield func(constraint.Expr) bool) {
		for cmt := range comments.IterNode(node, buildTagCommentFilter) {
			c, err := constraint.Parse(cmt)
			if err != nil {
				panic(err)
			}

			if !filter(c) {
				continue
			}

			if !yield(c) {
				return
			}
		}
	}
}

func buildTagCommentFilter(line string) bool {
	return strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build")
}
