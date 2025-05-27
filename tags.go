package evil

import (
	"go/ast"
	"go/build/constraint"
	"strings"
)

// parseTagsString returns a list of build tags from a comma separated or space separated string
func parseTagString(s string) []string {
	var tags []string

	// satisfy both comma separated and space separated conditions as
	// specified in the go build constraints: https://pkg.go.dev/cmd/go#hdr-Build_constraints
	for _, elem := range strings.Split(s, ",") {
		elem = strings.TrimSpace(elem)
		if elem != "" {
			for _, lit := range strings.Split(elem, " ") {
				tags = append(tags, strings.TrimSpace(lit))
			}
		}
	}

	return tags
}

// MatchTagsInFileNode returns true if the file matches any of the given build tags
func MatchTagsInFileNode(node *ast.File, tags []string, ignore []string) (bool, error) {
	if tags == nil && ignore == nil {
		// no build tags to check, and no build tags to ignore
		return true, nil
	}

	var exprs []constraint.Expr
	for _, cg := range node.Comments {
		for _, cmt := range cg.List {
			txt := cmt.Text

			if !commentIsBuildTag(txt) {
				continue
			}

			cst, err := constraint.Parse(txt)
			if err != nil {
				return false, err
			}

			exprs = append(exprs, cst)
		}
	}

	if len(exprs) == 0 {
		// no build tags to check
		// return true if no tags were specified
		return tags == nil, nil
	}

	return !tagsMatchAnyExpr(ignore, exprs) && tagsMatchAllExpr(tags, exprs), nil
}

func commentIsBuildTag(line string) bool {
	return strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build")
}

func tagsMatchAnyExpr(tags []string, constraints []constraint.Expr) bool {
	for _, v := range tags {
		for _, con := range constraints {
			if con.Eval(func(t string) bool {
				return t == v
			}) {
				return true
			}
		}
	}
	return false
}

func tagsMatchAllExpr(tags []string, constraints []constraint.Expr) bool {
	if len(tags) == 0 || len(constraints) == 0 {
		// no tags to match or file is unconstrained
		return true
	}

	cnt := 0
	for _, c := range constraints {
		for _, v := range tags {
			if c.Eval(func(t string) bool {
				return t == v
			}) {
				cnt++
				break
			}
		}
	}

	return cnt == len(constraints)
}
