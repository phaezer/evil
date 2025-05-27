package evil

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Generator is initialized in the temp main func and passed to the
//
//	calls to the Init and Register functions
type Generator struct {
	// fs is the file set used to parse the files
	fs *token.FileSet
	// buf is a buffer that holds the generated code
	buf strings.Builder
	// files is a map of package files to their parsed asts
	files []*ast.File
	// pkg is the name of the package that generated files will use
	pkg string
	// filename is the name of the path that generated file will be written to
	filename string
	// caller holds the info of the function that instantiated the Generator
	caller *callerInfo
}

func OutputName(name string) func(*Generator) {
	return func(g *Generator) {
		g.filename = name
	}
}

func PackageName(name string) func(*Generator) {
	return func(g *Generator) {
		g.pkg = name
	}
}

// NewGenerator creates a new Generator
func NewGenerator(opts ...func(*Generator)) *Generator {
	g := &Generator{
		fs:     token.NewFileSet(),
		caller: getCallerInfo(1),
	}

	for _, opt := range opts {
		opt(g)
	}

	if g.filename == "" {
		// if output was not provided check if GOEVIL_OUTPUT was set,
		// if not then use the caller's filename without .go extension + "__generated.go"

		if out := os.Getenv("GOEVIL_OUTPUT"); out != "" {
			g.filename = out
		} else {
			// use the caller's filename without the .go extension + "__generated.go"
			g.filename = fmt.Sprintf("%s__generated.go",
				strings.TrimSuffix(g.caller.fn.name, ".go"))
		}
	}

	if g.pkg == "" {
		// if package was not provided, use GOEVIL_PACKAGE env var if set
		// this should be set if go generate was used
		if out := os.Getenv("GOEVIL_PACKAGE"); out != "" {
			g.pkg = out
		} else {
			// assume it's the main package if not specified
			g.pkg = "main"
		}
	}

	return g
}

//
//// InitVar is used to initialize a variable in the global scope
////
////go:noinline
//func InitVar[T any](g *Generator, set T, to T) {
//	pc, file, line, ok := runtime.Caller(1)
//	if !ok {
//		panic("could not get calling function's PC")
//	}
//
//	log.Printf("pc = %#v, file = %#v\n, line = %#v", pc, file, line)
//
//	f := g.loadFile(file)
//
//	fnName, err := funcNameForPC(pc, f.ast.Name.Name)
//	if err != nil {
//		panic(err)
//	}
//
//	// find the caller function in the target file
//	callerFn := f.FindFnWithName(fnName)
//	if callerFn == nil {
//		panic(fmt.Sprintf("could not find caller function %s in target file", fnName))
//	}
//
//	matchingFn := f.MatchingCallExpr(fnName, line, callerFn)
//
//	if matchingFn == nil {
//		panic(fmt.Sprintf("could not find matching fn call expression in file %s on line %d", file, line))
//	}
//}

func funcNameForPC(pc uintptr, filePkgName string) (string, error) {
	callerFnName := runtime.FuncForPC(pc).Name()
	if len(callerFnName)-len(filePkgName) < 2 {
		return "", fmt.Errorf("function package name mismatch: %s != %s", callerFnName, filePkgName)
	}

	lastDot := strings.LastIndexByte(callerFnName, '.')
	fnName := callerFnName[lastDot+1:]

	// make sure that the function name starts with the package name
	if !strings.HasSuffix(callerFnName, filePkgName+"."+fnName) {
		return "", fmt.Errorf("function package name mismatch: %s != %s", callerFnName, filePkgName)
	}

	return fnName, nil
}

// getCurrentPackagePath returns the full import path of the package containing the calling code
func getBuildPackageFromFile(file string) (*build.Package, error) {
	// Get caller information
	// Convert to absolute path if it's not already
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %w", err)
	}

	// Find the directory containing the file
	dir := filepath.Dir(absPath)

	// Use go/build to find the package import path
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		return nil, fmt.Errorf("could not determine package path: %w", err)
	}

	return pkg, nil
}
