package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"iter"

	"github.com/phaezer/evil/internal/parser/tags"
	"github.com/phaezer/evil/internal/seq"
)

// WithBuildTags returns a NodeFilter that returns true if the file has any of the given build tags
func WithBuildTags(v ...string) func(*ast.File) bool {
	return func(file *ast.File) bool {
		for range tags.IterInFileNode(file, tags.WithAny(v...)) {
			// if there are any results from the filter, return true
			return true
		}

		return false
	}
}

func IterFileNodesInPath(
	fs *token.FileSet,
	path string,
	mode parser.Mode,
	fileFilter func(string) bool,
	nodeFilter func(*ast.File) bool,
) iter.Seq2[*ast.File, error] {
	if isDir(path) {
		return iterFileNodesInDir(fs, path, mode, fileFilter, nodeFilter)
	}

	// if the path is a file, parse it
	file, err := ParseFileWithFilter(fs, path, nil, mode, nodeFilter)

	// return the single file wrapped in a sequence
	return seq.Once2(file, err)
}

func iterFileNodesInDir(
	fs *token.FileSet,
	dir string,
	mode parser.Mode,
	fileFilter func(string) bool,
	nodeFilter func(*ast.File) bool,
) iter.Seq2[*ast.File, error] {
	return func(yield func(*ast.File, error) bool) {
		// we're assuming dir is a directory
		for f := range iterGoFilesInDir(dir, fileFilter) {
			file, err := ParseFileWithFilter(fs, f, nil, mode, nodeFilter)
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			if file != nil && !yield(file, nil) {
				return
			}
		}
	}
}

func ParseFileWithFilter(
	fs *token.FileSet,
	filename string,
	src any,
	mode parser.Mode,
	filter func(*ast.File) bool,
) (*ast.File, error) {
	// test if the file passes the filter
	ok, err := testFile(filename, src, mode, filter)
	if err != nil || !ok {
		return nil, err
	}

	return parser.ParseFile(fs, filename, src, mode)
}

func testFile(filename string, src any, mode parser.Mode, filter func(*ast.File) bool) (bool, error) {
	if filter == nil {
		// no filter passed; no need to parse to check the filter
		return true, nil
	}

	// use a temporary fileset to parse the token before checking parsed file with the filter
	file, err := parser.ParseFile(token.NewFileSet(), filename, src, mode)
	if err != nil {
		// parser error
		return false, err
	}

	if filter(file) {
		return true, nil
	}

	return false, nil
}
