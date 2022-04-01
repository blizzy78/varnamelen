package varnamelen

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// constant represents a declared constant.
type constant struct {
	// name is the name of the constant.
	name string

	// valueSpec is the value specification that declares the constant.
	valueSpec *ast.ValueSpec
}

// checkConstants applies the analysis to constants in constantToDist, according to cfg.
func checkConstants(pass *analysis.Pass, constantToDist map[constant]int, cfg configuration) {
	for cons, dist := range constantToDist {
		if cfg.ignoreNames.contains(cons.name) {
			continue
		}

		if cfg.ignoreDecls.matchConstant(cons) {
			continue
		}

		if checkNameAndDistance(cons.name, dist, cfg) {
			continue
		}

		pass.Reportf(cons.valueSpec.Pos(), "constant name '%s' is too short for the scope of its usage", cons.name)
	}
}

// match returns true if c matches decl.
func (c constant) match(decl identDeclaration) bool {
	if c.name != decl.name {
		return false
	}

	if !decl.constant {
		return false
	}

	return true
}
