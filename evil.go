package evil

import (
	"fmt"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const (
	packageName = "evil"
	packagePath = "github.com/phaezer/" + packageName
)

type Evil struct {
	fs *token.FileSet
	// buf is a buffer that holds the generated code
	buf strings.Builder
	// pkg is the name of the package that generated files will use
	pkg string
	// tags is a list of build tags to include
	tags []string

	// files is a map of go files that have been parsed
	files map[string]*File
	// filename is the path the generated file will be written to
	filename string

	caller *callerInfo
	// build holds the [debug.BuildInfo] of the running binary
	build *debug.BuildInfo
}

func New(filename string) *Evil {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic(fmt.Errorf("could not read build info from running binary"))
	}

	evl := &Evil{
		fs:    token.NewFileSet(),
		files: map[string]*File{},
		build: bi,
	}

	if len(evl.tags) == 0 {
		// include the tags used to execute this file
		for _, bs := range bi.Settings {
			if bs.Key == "-tags" {
				tags := strings.Split(bs.Value, " ")
				for _, tag := range tags {
					evl.tags = append(evl.tags, strings.TrimSpace(tag))
				}
			}
		}
	}

	return evl
}

//
//func (e *Evil) buildTagFilter() func(expr constraint.Expr) bool {
//	// we should always have at least one excluded tag
//	negativeFilter := func(expr constraint.Expr) bool {
//		for _, tag := range e.excludeTags {
//			if expr.Eval(func(t string) bool {
//				return t == tag
//			}) {
//				return false
//			}
//		}
//		return true
//	}
//
//	// if we we don't have any tags, we can use the negative filter
//	if len(e.tags) == 0 {
//		return negativeFilter
//	}
//
//	positiveFilter := func(expr constraint.Expr) bool {
//		for _, tag := range e.tags {
//			if expr.Eval(func(t string) bool {
//				return t == tag
//			}) {
//				return true
//			}
//		}
//		return false
//	}
//
//	// if we have both tags and exclude tags, we can use both filters
//	return func(expr constraint.Expr) bool {
//		return positiveFilter(expr) && negativeFilter(expr)
//	}
//}
//
//func (e *Evil) ParseDir(dir string) error {
//	// always use the absolute path
//	abs, err := filepath.Abs(dir)
//	if err != nil {
//		return err
//	}
//
//	buildTagFilter := e.buildTagFilter()
//
//	// iter over the go files in the directory
//	err = inspectDir(abs, func(absPath string, entry os.DirEntry) bool {
//		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
//			return false
//		}
//
//		var f *File
//		f, err = ParseFile(e.fs, absPath, nil, parser.ParseComments)
//		if err != nil {
//			return false
//		}
//
//		// collect build tags in the file
//		buildTags := []string{}
//		for line := range IterComments(f.node, commentIsBuildTag) {
//			buildTags = append(buildTags, line)
//		}
//
//		// separate build tags in the same file are AND-ed
//		// so if any tag matches the filter, include the file
//		for _, tag := range buildTags {
//			c, err := constraint.Parse(tag)
//			if err != nil {
//				// invalid build tag
//				panic(err)
//			}
//
//			// if the build tag doesn't match the filter, remove the file
//			if !buildTagFilter(c) {
//				e.fs.RemoveFile(f.token)
//				return false
//			}
//		}
//
//		// file passed the filter, add the file
//		e.files[absPath] = f
//		return true
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// inspectDir iterates over the files in a directory and calls fn for each file
// if fn returns false, the iteration stops
// fn is called with the absolute path to the file and the [os.DirEntry]
//
// from: https://gist.github.com/phaezer/077d31c09d2d8d378f3541fa2293ee04
func inspectDir(dir string, fn func(path string, entry os.DirEntry) bool) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	stat, err := os.Stat(abs)
	if err != nil {
		// unlikely to have a [*os.PathError] if [filepath.Abs] didn't return an error
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		p := path.Join(abs, entry.Name())
		if !fn(p, entry) {
			return nil
		}
	}

	return nil
}
