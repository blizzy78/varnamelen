package varnamelen

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// variable represents a declared variable.
type variable struct {
	// name is the name of the variable.
	name string

	// typ is the type of the variable.
	typ string

	// assign is the assign statement that declares the variable.
	assign *ast.AssignStmt

	// valueSpec is the value specification that declares the variable.
	valueSpec *ast.ValueSpec
}

// checkVariables applies the analysis to variables in varToDist, according to cfg.
func checkVariables(pass *analysis.Pass, varToDist map[variable]int, cfg configuration) { //nolint:gocognit // it's not that complicated
	for variable, dist := range varToDist {
		if variable.isConventional() {
			continue
		}

		if cfg.ignoreNames.contains(variable.name) {
			continue
		}

		if cfg.ignoreDecls.matchVariable(variable) {
			continue
		}

		if checkTypeAssertOk(variable, cfg) {
			continue
		}

		if checkMapIndexOk(variable, cfg) {
			continue
		}

		if checkChannelReceiveOk(variable, cfg) {
			continue
		}

		if checkNameAndDistance(variable.name, dist, cfg) {
			continue
		}

		var pos token.Pos
		if variable.assign != nil {
			pos = variable.assign.Pos()
		} else {
			pos = variable.valueSpec.Pos()
		}

		pass.Reportf(pos, "variable name '%s' is too short for the scope of its usage", variable.name)
	}
}

// checkTypeAssertOk returns true if "ok" variables that hold the bool return value of a type assertion
// should be ignored, and if vari is such a variable.
func checkTypeAssertOk(vari variable, cfg configuration) bool {
	return cfg.ignoreTypeAssertOk && vari.isTypeAssertOk()
}

// checkMapIndexOk returns true if "ok" variables that hold the bool return value of a map index
// should be ignored, and if vari is such a variable.
func checkMapIndexOk(vari variable, cfg configuration) bool {
	return cfg.ignoreMapIndexOk && vari.isMapIndexOk()
}

// checkChannelReceiveOk returns true if "ok" variables that hold the bool return value of a channel receive
// should be ignored, and if vari is such a variable.
func checkChannelReceiveOk(vari variable, cfg configuration) bool {
	return cfg.ignoreChanRecvOk && vari.isChannelReceiveOk()
}

// match returns true if v matches decl.
func (v variable) match(decl identDeclaration) bool {
	if v.name != decl.name {
		return false
	}

	if decl.constant {
		return false
	}

	if v.typ == "" {
		return false
	}

	return decl.matchType(v.typ)
}

// isTypeAssertOk returns true if v is an "ok" variable that holds the bool return value of a type assertion.
func (v variable) isTypeAssertOk() bool {
	if v.name != "ok" {
		return false
	}

	if v.assign == nil {
		return false
	}

	if len(v.assign.Lhs) != 2 {
		return false
	}

	ident, ok := v.assign.Lhs[1].(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "ok" {
		return false
	}

	if len(v.assign.Rhs) != 1 {
		return false
	}

	if _, ok := v.assign.Rhs[0].(*ast.TypeAssertExpr); !ok {
		return false
	}

	return true
}

// isMapIndexOk returns true if v is an "ok" variable that holds the bool return value of a map index.
func (v variable) isMapIndexOk() bool {
	if v.name != "ok" {
		return false
	}

	if v.assign == nil {
		return false
	}

	if len(v.assign.Lhs) != 2 {
		return false
	}

	ident, ok := v.assign.Lhs[1].(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "ok" {
		return false
	}

	if len(v.assign.Rhs) != 1 {
		return false
	}

	if _, ok := v.assign.Rhs[0].(*ast.IndexExpr); !ok {
		return false
	}

	return true
}

// isChannelReceiveOk returns true if v is an "ok" variable that holds the bool return value of a channel receive.
func (v variable) isChannelReceiveOk() bool {
	if v.name != "ok" {
		return false
	}

	if v.assign == nil {
		return false
	}

	if len(v.assign.Lhs) != 2 {
		return false
	}

	ident, ok := v.assign.Lhs[1].(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "ok" {
		return false
	}

	if len(v.assign.Rhs) != 1 {
		return false
	}

	unary, ok := v.assign.Rhs[0].(*ast.UnaryExpr)
	if !ok {
		return false
	}

	if unary.Op != token.ARROW {
		return false
	}

	return true
}

// isConventional returns true if v matches a conventional Go variable name and type,
// such as "ctx context.Context" or "t *testing.T".
func (v variable) isConventional() bool {
	for _, decl := range conventionalDecls {
		if v.match(decl) {
			return true
		}
	}

	return false
}
