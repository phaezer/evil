package evil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// File represents a parsed Go file
// this wraps the [ast.File] for the file and includes the [token.File] struct
type File struct {
	// node is the [ast.File] for the file
	node *ast.File
	// token is the [token.File] fileset file for as registered with [parser.ParseFile]
	token *token.File
}

// ParseFile parses a file with a given filename or src and parses it into a [File]
//
//	and adds the token to the [token.FileSet]
func ParseFile(fs *token.FileSet, filename string, src any, mode parser.Mode) (*File, error) {
	fileAst, err := parser.ParseFile(fs, filename, src, mode)
	if err != nil {
		return nil, fmt.Errorf("parsing token %s failed: %w", filename, err)
	}

	f := &File{
		node: fileAst,
	}

	// get the ref to the parsed File in the token set with the filename
	fs.Iterate(func(file *token.File) bool {
		if file.Name() == filename {
			f.token = file
			return false
		}
		return true
	})

	if f.token == nil {
		panic(fmt.Sprintf("failed finding token for file %s", f.Name()))
	}

	return f, nil
}

// Inspect calls [ast.Inspect] on the file with the given function
func (f *File) Inspect(fn func(n ast.Node) bool) {
	ast.Inspect(f.node, fn)
}

// Name returns the file name of file f as passed to [parser.ParseFile].
func (f *File) Name() string {
	return f.token.Name()
}

// PackageName returns the package name of the file as defined in the node
func (f *File) PackageName() string {
	return f.node.Name.String()
}
