package fnode

import (
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"os"
	"path"
	"path/filepath"

	"github.com/phaezer/evil/tools/filter"
	"github.com/phaezer/evil/tools/parse/constraints"
)

type NodeFilter filter.Filter[*ast.File]

func WithBuildTags(tags []string) NodeFilter {
	return func(file *ast.File) bool {
		for range constraints.IterBuildTags(file, constraints.WithAnyBuildTags(tags)) {
			// if there are any results from the filter, return true
			return true
		}

		return false
	}
}

func IterNodesInPath(
	fs *token.FileSet,
	path string,
	mode parser.Mode,
	fileFilter filter.Filter[string],
	nodeFilter NodeFilter,
) iter.Seq[*ast.File] {
	// if nodeFilter is nil, all files are returned
	if nodeFilter == nil {
		nodeFilter = func(*ast.File) bool { return true }
	}

	if isDir(path) {
		return iterNodesInDir(fs, path, mode, fileFilter, nodeFilter)
	}

	// if the path is a file, parse it
	file, err := parser.ParseFile(fs, path, nil, mode)
	if err != nil {
		panic(err)
	}

	// yield the file
	return func(yield func(*ast.File) bool) {
		if nodeFilter(file) {
			if !yield(file) {
				return
			}
		}
	}
}

func iterNodesInDir(
	fs *token.FileSet,
	dir string,
	mode parser.Mode,
	fileFilter filter.Filter[string],
	nodeFilter NodeFilter,
) iter.Seq[*ast.File] {
	return func(yield func(*ast.File) bool) {
		// we're assuming dir is a directory
		for f := range iterGoFilesInDir(dir, fileFilter) {
			file, err := parseFileWithFilter(fs, f, nil, mode, nodeFilter)
			if err != nil {
				panic(err)
			}

			if file == nil {
				continue
			}

			if !yield(file) {
				return
			}
		}
	}
}

func parseFileWithFilter(fs *token.FileSet, filename string, src any, mode parser.Mode, filter NodeFilter) (*ast.File, error) {
	// use a temporary fileset to parse the file
	// only the files that pass the filter are added to the original fileset
	file, err := parser.ParseFile(token.NewFileSet(), filename, src, mode)
	if err != nil {
		return nil, err
	}

	if filter == nil || filter(file) {
		// parse the file to the original fileset
		file, err = parser.ParseFile(fs, filename, src, mode)
		if err != nil {
			// should never happen if original fileset parser was successful
			// but errors should not be ignored
			panic(err)
		}

		return file, nil
	}

	return nil, nil

}

func iterGoFilesInDir(dir string, filter filter.Filter[string]) iter.Seq[string] {
	// we're assuming dir is a directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// if filter is nil, all files are returned
	if filter == nil {
		filter = func(string) bool { return true }
	}

	return func(yield func(string) bool) {
		for _, f := range entries {
			if f.IsDir() || filepath.Ext(f.Name()) != ".go" {
				continue
			}

			p := path.Join(dir, f.Name())

			if filter(p) {
				if !yield(p) {
					return
				}
			}
		}
	}
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	// todo: should we panic instead of returning false?
	return err == nil && s.IsDir()
}
