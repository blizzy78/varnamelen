package varnamelen

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// receiver represents a declared method receiver.
type receiver struct {
	// name is the name of the receiver.
	name string

	// typ is the type of the receiver.
	typ string

	// field is the declaration of the receiver.
	field *ast.Field
}

// checkReceivers applies the analysis to method receivers in receiverToDist, according to cfg.
func checkReceivers(pass *analysis.Pass, receiverToDist map[receiver]int, cfg configuration) {
	for recv, dist := range receiverToDist {
		if cfg.ignoreNames.contains(recv.name) {
			continue
		}

		if cfg.ignoreDecls.matchReceiver(recv) {
			continue
		}

		if checkNameAndDistance(recv.name, dist, cfg) {
			continue
		}

		pass.Reportf(recv.field.Pos(), "method receiver name '%s' is too short for the scope of its usage", recv.name)
	}
}

// match returns whether r matches decl.
func (r receiver) match(decl identDeclaration) bool {
	if r.name != decl.name {
		return false
	}

	return decl.matchType(r.typ)
}
