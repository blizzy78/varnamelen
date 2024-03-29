package varnamelen

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// typeParam represents a declared type parameter.
type typeParam struct {
	// name is the name of the type parameter.
	name string

	// typ is the type of the type parameter.
	typ string

	// field is the field that declares the type parameter.
	field *ast.Field
}

// match returns whether p matches decl.
func (p typeParam) match(decl identDeclaration) bool {
	if p.name != decl.name {
		return false
	}

	return decl.matchType(p.typ)
}

// checkTypeParams applies the analysis to type parameters in paramToDist, according to cfg.
func checkTypeParams(pass *analysis.Pass, paramToDist map[typeParam]int, cfg configuration) {
	for param, dist := range paramToDist {
		if cfg.ignoreNames.contains(param.name) {
			continue
		}

		if cfg.ignoreDecls.matchTypeParameter(param) {
			continue
		}

		if checkNameAndDistance(param.name, dist, cfg) {
			continue
		}

		pass.Reportf(param.field.Pos(), "type parameter name '%s' is too short for the scope of its usage", param.name)
	}
}

// isTypeParam returns true if field is a type parameter of any of the given funcs.
func isTypeParam(field *ast.Field, funcs []*ast.FuncDecl, funcLits []*ast.FuncLit) bool { //nolint:gocognit // it's not that complicated
	for _, f := range funcs {
		if f.Type.TypeParams == nil {
			continue
		}

		for _, p := range f.Type.TypeParams.List {
			if p == field {
				return true
			}
		}
	}

	for _, f := range funcLits {
		if f.Type.TypeParams == nil {
			continue
		}

		for _, p := range f.Type.TypeParams.List {
			if p == field {
				return true
			}
		}
	}

	return false
}
