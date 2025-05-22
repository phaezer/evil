package evil

import (
	"go/ast"
	"iter"
)

// todo: implement a Visitor instead of iterators
//  see: https://cs.opensource.google/go/go/+/refs/tags/go1.24.3:src/go/ast/walk.go;l=33

func iterAstDeclsForType[E ast.Decl](decls []ast.Decl) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, decl := range decls {
			if dt, ok := decl.(E); ok {
				if !yield(dt) {
					return
				}
			}
		}
	}
}

func iterAstFnNames(decls []ast.Decl) iter.Seq[string] {
	return func(yield func(name string) bool) {
		for fn := range iterAstDeclsForType[*ast.FuncDecl](decls) {
			if !yield(fn.Name.Name) {
				return
			}
		}
	}
}

func iterAstStmtFields[E any](stmt ast.Stmt, filter func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		// loop through all stmts in the body list recursively
		// all exprs nested in each stmt should be yielded from [iterStructFieldsForInterface]
		for expr := range iterStructFieldsForInterface[ast.Expr](stmt) {
			if v, ok := expr.(E); ok {
				if filter == nil || filter(v) {
					if !yield(v) {
						return
					}
				}
			}
		}
	}
}
