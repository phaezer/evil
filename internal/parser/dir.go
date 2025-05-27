package parser

import (
	"go/ast"
	"iter"
	"os"
	"path"
	"path/filepath"
)

func iterGoFilesInDir(dir string, f func(string) bool) iter.Seq[string] {
	// we're assuming dir is a directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	return func(yield func(string) bool) {
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
				continue
			}

			p := path.Join(dir, entry.Name())

			if f == nil || f(p) {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// GetImportName returns the import name for the given package path
func GetImportName(file *ast.File, pkgPath string) string {
	if file == nil || file.Imports == nil || len(file.Imports) == 0 {
		return ""
	}

	for _, imp := range file.Imports {
		if imp.Path.Value == `"`+pkgPath+`"` {
			if imp.Name != nil {
				// an import alias was used, we'll use this instead when we walk the ast
				return imp.Name.Name
			}
			return pkgPath
		}
	}

	return ""
}
