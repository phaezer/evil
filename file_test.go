package evil

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	fs *token.FileSet
}

func (suite *FileTestSuite) SetupTest() {
	// get the abs path of the token
	suite.fs = token.NewFileSet()
}

//		log.Printf("package name: %#v\n", suite.file.node.Name.Name)
//		log.Printf("import path: %#v\n", suite.file.filename)
//		log.Printf("found functions: %v\n", slices.Collect(suite.file.IterFnNames()))
//	}
func (suite *FileTestSuite) TestParseFile() {
	t := suite.T()

	absPath, err := filepath.Abs("examples/simple/fib.go")
	assert.NoError(t, err)

	f, err := ParseFile(suite.fs, absPath, nil, parser.ParseComments)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.NotNil(t, f.node)
	assert.NotNil(t, f.token)
	assert.Equal(t, absPath, f.Name())
}

func (suite *FileTestSuite) TestParseDir() {
	t := suite.T()

	dir, err := filepath.Abs("examples/simple")
	assert.NoError(t, err)

	files, err := ParseDir(suite.fs, dir, nil)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestFileSuiteSimple(t *testing.T) {
	s := &FileTestSuite{}

	suite.Run(t, s)
}
