package evil

import (
	"errors"
	"go/token"
	"strings"
)

const (
	packageName = "evil"
	packagePath = "github.com/phaezer/" + packageName
)

// IgnoreTags is a list of build tags to ignore when parsing go files
var IgnoreTags = []string{"ignore", packageName}

type Evil struct {
	fs *token.FileSet
	// buf is a buffer that holds the generated code
	buf strings.Builder
	// pkg is the name of the package that generated files will use
	pkg string
	// tags is a list of build tags to include when parsing go files
	tags []string
	// files is a map of go files that have been parsed
	files map[string]*File
	// filename is the path the generated file will be written to
	filename string

	// caller *callerInfo
	// build holds the [debug.BuildInfo] of the running binary
	// build *debug.BuildInfo
}

func New(filename string, tags []string) *Evil {
	evl := &Evil{
		fs:    token.NewFileSet(),
		files: map[string]*File{},
		tags:  tags,
	}

	return evl
}

var readDebugBuildError = errors.New("failed reading build info from running binary")
