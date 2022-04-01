package varnamelen

import (
	"strings"
)

// identDeclaration is an ident declaration.
type identDeclaration struct {
	// name is the name of the ident.
	name string

	// constant is true if the declaration is actually declaring a constant.
	constant bool

	// typ is the type of the ident. Not used for constants.
	typ string
}

// A declarationError is returned when an ident declaration cannot be parsed.
type declarationError string

// mustParseIdentDeclaration parses and returns an ident declaration parsed from decl.
func mustParseIdentDeclaration(decl string) identDeclaration {
	d, _ := parseIdentDeclaration(decl)
	return d
}

// parseIdentDeclaration parses and returns an ident declaration parsed from decl.
func parseIdentDeclaration(decl string) (identDeclaration, error) {
	if strings.HasPrefix(decl, "const ") {
		name := strings.TrimPrefix(decl, "const ")
		if strings.TrimSpace(name) == "" {
			return identDeclaration{}, declarationError(decl)
		}

		return identDeclaration{
			name:     name,
			constant: true,
		}, nil
	}

	parts := strings.SplitN(decl, " ", 2)

	if len(parts) != 2 {
		return identDeclaration{}, declarationError(decl)
	}

	return identDeclaration{
		name: parts[0],
		typ:  parts[1],
	}, nil
}

// matchType returns true if typ matches d.typ.
func (d identDeclaration) matchType(typ string) bool {
	return d.typ == typ
}

// Error implements error.
func (e declarationError) Error() string {
	return string(e) + ": invalid declaration"
}
