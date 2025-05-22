package evil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"path/filepath"
	"sync"
)

// File is a parsed Go file with associated metadata
// File includes the [ast.File] for the file
type File struct {
	mu sync.Mutex
	// filename is the filename of the file used when [parser.ParseFile] is calleds
	filename string
	// ast is the ast.File for the file
	ast *ast.File
	// file is the token.File for the file in the file set
	file *token.File
}

// ImportName returns the import name used in the file for this package
func (f *File) ImportName() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return getImportNameFromFileAst(f.ast)
}

func ParseFile(fs *token.FileSet, filename string) (*File, error) {
	// get absolute path of the file
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	// parse the file with comments
	fileAst, err := parser.ParseFile(fs, absPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file %s failed: %w", absPath, err)
	}

	f := &File{
		filename: absPath,
		ast:      fileAst,
	}

	// get the ref to the parsed File in the file set with the filename
	fs.Iterate(func(file *token.File) bool {
		if file.Name() == f.filename {
			f.file = file
			return false
		}
		return true
	})

	if f.ast == nil {
		return nil, fmt.Errorf("could not parse file %s", absPath)
	}

	if f.ast.Decls == nil || len(f.ast.Decls) == 0 {
		return nil, fmt.Errorf("no decls in file %s", absPath)
	}

	return f, nil
}

func (f *File) FindFnWithName(fnName string) *ast.FuncDecl {
	for fn := range iterAstDeclsForType[*ast.FuncDecl](f.ast.Decls) {
		if fn.Name.Name == fnName {
			return fn
		}
	}
	return nil
}

//
//func (f *File) MatchingCallExpr(fnName string, line int, callerFn *ast.FuncDecl) *ast.CallExpr {
//	for _, stmt := range callerFn.Body.List {
//		for expr := range iterStructFieldsForInterface[ast.Expr](stmt) {
//			if callExpr, ok := expr.(*ast.CallExpr); ok {
//				fmt.Printf("callExpr = %#v\n", callExpr)
//				if selExp, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
//					fnLine := f.file.Line(callExpr.Pos())
//					if selExp.Name == fnName && fnLine == line {
//						return callExpr
//					}
//				}
//			}
//		}
//	}
//	return nil
//}

func (f *File) MatchingCallExpr(fnName string, line int, callerFn *ast.FuncDecl) *ast.CallExpr {
	for _, stmt := range callerFn.Body.List {
		for expr := range iterStructFieldsForInterface[ast.Expr](stmt) {
			if callExpr, ok := expr.(*ast.CallExpr); ok {
				fmt.Printf("callExpr = %#v\n", callExpr)
				if selExp, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					fnLine := f.file.Line(callExpr.Pos())
					if selExp.Sel.Name == fnName && fnLine == line {
						return callExpr
					}
				}
			}
		}
	}
	return nil
}

func (f *File) IterFnNames() iter.Seq[string] {
	return iterAstFnNames(f.ast.Decls)
}

func getImportNameFromFileAst(file *ast.File) string {
	if file == nil || file.Imports == nil || len(file.Imports) == 0 {
		return ""
	}

	for _, imp := range file.Imports {
		if imp.Path.Value == `"`+packagePath+`"` {
			if imp.Name != nil {
				// an import alias was used, we'll use this instead when we walk the ast
				return imp.Name.Name
			}
			return importName
		}
	}

	return ""
}
