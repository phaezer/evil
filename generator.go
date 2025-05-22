package evil

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"iter"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Generator is initialized in the temp main func and passed to the
//
//	calls to the Init and Register functions
type Generator struct {
	mu sync.RWMutex
	fs *token.FileSet

	buf strings.Builder

	// files is a map of package files to their parsed asts
	files map[string]*File

	pkgName string
	output  string
}

// NewGenerator creates a new Generator
func NewGenerator(outputTo string, pkgName string) *Generator {
	g := &Generator{
		fs:      token.NewFileSet(),
		files:   make(map[string]*File),
		pkgName: pkgName,
		output:  outputTo,
	}

	return g
}

func PackageFrom(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(abs)
	if err != nil {
		panic(err)
	}

	if stat.IsDir() {

	}

	// if the path is a directory, scan it for go files
	if stat, err := os.Stat(abs); err == nil && stat.IsDir() {
		// TODO: scan the directory for go files
	}
}

func PackageNameFromFile(file *ast.File) string {
	return file.Name.Name
}

// loadFile parses a file and adds it to the file map
func (g *Generator) loadFile(filename string) *File {
	g.mu.Lock()
	defer g.mu.Unlock()

	// check if the file has already been parsed; if it was, return it
	if f, ok := g.files[filename]; ok {
		return f
	}

	f, err := ParseFile(g.fs, filename)
	if err != nil {
		panic(fmt.Sprintf("could not parse file %s: %v", filename, err))
	}

	g.files[filename] = f
	return f
}

// InitVar is used to initialize a variable in the global scope
//
//go:noinline
func InitVar[T any](g *Generator, set T, to T) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("could not get calling function's PC")
	}

	log.Printf("pc = %#v, file = %#v\n, line = %#v", pc, file, line)

	f := g.loadFile(file)

	fnName, err := funcNameForPC(pc, f.ast.Name.Name)
	if err != nil {
		panic(err)
	}

	// find the caller function in the target file
	callerFn := f.FindFnWithName(fnName)
	if callerFn == nil {
		panic(fmt.Sprintf("could not find caller function %s in target file", fnName))
	}

	matchingFn := f.MatchingCallExpr(fnName, line, callerFn)

	if matchingFn == nil {
		panic(fmt.Sprintf("could not find matching fn call expression in file %s on line %d", file, line))
	}
}

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
