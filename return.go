package varnamelen

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// namedReturn represents a declared named return value.
type namedReturn struct {
	// name is the name of the return value.
	name string

	// typ is the type of the return value.
	typ string

	// field is the declaration of the return value.
	field *ast.Field
}

// checkReceivers applies the analysis to named return values in returnToDist, according to cfg.
func checkReturns(pass *analysis.Pass, returnToDist map[namedReturn]int, cfg configuration) {
	for ret, dist := range returnToDist {
		if cfg.ignoreNames.contains(ret.name) {
			continue
		}

		if cfg.ignoreDecls.matchNamedReturn(ret) {
			continue
		}

		if checkNameAndDistance(ret.name, dist, cfg) {
			continue
		}

		pass.Reportf(ret.field.Pos(), "return value name '%s' is too short for the scope of its usage", ret.name)
	}
}

// match returns whether r matches decl.
func (r namedReturn) match(decl identDeclaration) bool {
	if r.name != decl.name {
		return false
	}

	return decl.matchType(r.typ)
}
