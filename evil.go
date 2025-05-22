package evil

import (
	"go/ast"
	"reflect"
)

const (
	importName  = "evil"
	packagePath = "github.com/phaezer/" + importName
)

//type Generator interface {
//	Init(any, any)
//	Register(string, any)
//}

type initTarget[T any] struct {
	// name is the name of the variable
	name string
	// inputType is the reflect type of the variable to be set
	inputType reflect.Type
	// rVal is the reflect value will be set after calling the associated fn
	rVal reflect.Value
	// var ident that matches this variable
	ident *ast.Ident
}
