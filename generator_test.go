package evil

import (
	"log"
	"testing"
)

type (
	node interface {
		node()
	}

	nestedIface interface {
		nested()
	}

	child struct {
		N node
	}

	childSlice struct {
		N []node
	}

	terminalNode struct{}

	parentNode struct {
		TN *terminalNode
		CN *childNestedNode
	}

	childNestedNode struct{}
)

func (c *child) node()             {}
func (c *childSlice) node()        {}
func (c *terminalNode) node()      {}
func (c *childNestedNode) node()   {}
func (c *childNestedNode) nested() {}

func TestReflectStructFieldIterator(t *testing.T) {
	nodeWithTargets := &child{
		N: &child{
			N: &terminalNode{},
		},
	}

	nodeWithNestedTarget := &parentNode{
		TN: &terminalNode{},
		CN: &childNestedNode{},
	}

	for v := range iterStructFieldsForInterface[node](nodeWithTargets) {
		log.Printf("v = %#v\n", v)
	}

	for v := range iterStructFieldsForInterface[nestedIface](nodeWithNestedTarget) {
		log.Printf("nested v = %#v\n", v)
	}
}
