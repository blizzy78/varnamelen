package varnamelen

import (
	"go/ast"
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
