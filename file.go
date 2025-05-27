package evil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
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

// Inspect calls [ast.Inspect] on the [ast.File] node with the given fn
func (f *File) Inspect(fn func(n ast.Node) bool) {
	ast.Inspect(f.node, fn)
}

// Name returns the file name of file f as passed to [parser.ParseFile].
func (f *File) Name() string {
	return f.token.Name()
}

// PackageName returns the package name of the file as defined in the node
func (f *File) Package() string {
	return f.node.Name.String()
}

// Package represents a parsed Go package and its files
type Package struct {
	// name of the package as defined in the package declaration
	name string
	// path is the absolute filepath to the package
	path string
	// files is the list of files in the package
	files []*File
}

// ParseDir parses all of the go files in a directory with the given [token.FileSet]
// if tags is not nil, then only files that match any of the go:build tag constraints will be parsed
//
// ParseDir does not recurse into subdirectories
func ParseDir(fs *token.FileSet, path string, tags []string) ([]*Package, error) {
	// always use the absolute path
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// iterate over all of the go files in the directory
	dirIter, err := DirIterator(abs, func(path string, entry os.DirEntry) bool {
		return !entry.IsDir() && filepath.Ext(entry.Name()) == ".go"
	})
	if err != nil {
		return nil, err
	}

	pkgs := make(map[string]*Package)

	for goFile := range dirIter {
		f, err := ParseFile(fs, goFile, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		matched, err := MatchTagsInFileNode(f.node, tags, IgnoreTags)
		if err != nil {
			return nil, err
		}

		if !matched {
			// remove the file from the token set
			fs.RemoveFile(f.token)
			continue
		}

		if _, ok := pkgs[f.Package()]; !ok {
			pkgs[f.Package()] = &Package{
				name: f.Package(),
				path: path,
			}
		}

		pkgs[f.Package()].files = append(pkgs[f.Package()].files, f)
	}

	return slices.Collect(maps.Values(pkgs)), nil
}

// NotADirError is returned when a path is not a directory
type NotADirError struct {
	path string
}

func (e NotADirError) Error() string {
	return fmt.Sprintf("%s is not a directory", e.path)
}

// DirIterator iterates over the files in a directory and calls fn for each file
// if fn returns false, the iteration stops
// fn is called with the absolute path to the file and the [os.DirEntry]
//
// from: https://gist.github.com/phaezer/077d31c09d2d8d378f3541fa2293ee04
func DirIterator(dir string, f func(string, os.DirEntry) bool) (iter.Seq[string], error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(abs)
	if err != nil {
		// unlikely to have a [*os.PathError] if [filepath.Abs] didn't return an error
		return nil, err
	}

	if !stat.IsDir() {
		return nil, &NotADirError{path: dir}
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {
		for _, entry := range entries {
			p := path.Join(abs, entry.Name())
			if f == nil || f(p, entry) {
				if !yield(p) {
					return
				}
			}
		}
	}, nil
}
