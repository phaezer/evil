package evil

import (
	"go/token"

	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	FilePath string

	fs   *token.FileSet
	file *File
}

//
//func (suite *FileTestSuite) SetupTest() {
//	// get the abs path of the token
//	suite.fs = token.NewFileSet()
//
//	p, err := filepath.Abs(suite.FilePath)
//	if err != nil {
//		suite.T().Fatal(err)
//	}
//	log.Println("example path = ", p)
//
//	suite.file, err = ParseFile(suite.fs, p)
//	if err != nil {
//		suite.T().Fatal(err)
//	}
//
//	log.Printf("package name: %#v\n", suite.file.node.Name.Name)
//	log.Printf("import path: %#v\n", suite.file.filename)
//	log.Printf("found functions: %v\n", slices.Collect(suite.file.IterFnNames()))
//}
//
//func (suite *FileTestSuite) findFn(name string) *ast.FuncDecl {
//	fn := suite.file.FindFnWithName(name)
//	log.Printf("found fn with name = %#v\n", fn.Name.Name)
//	assert.NotNil(suite.T(), fn)
//	return fn
//}
//
//func (suite *FileTestSuite) TestFindFnWithName() {
//	_ = suite.findFn("evilSetVarFibonacci")
//}
//
//func (suite *FileTestSuite) TestFile_MatchingCallExpr() {
//	fn := suite.findFn("evilSetVarFibonacci")
//
//	matchingFn := suite.file.MatchingCallExpr("Init", 36, fn)
//
//	if matchingFn == nil {
//		suite.T().Fatal("could not find matching fn call")
//	}
//}
//
//func TestFileSuiteSimple(t *testing.T) {
//	s := &FileTestSuite{
//		FilePath: "examples/simple/fib.go",
//	}
//
//	suite.Run(t, s)
//}
