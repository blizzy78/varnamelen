package varnamelen

import (
	"go/ast"
	"go/types"
	"strings"
)

// importDeclaration is an import declaration.
type importDeclaration struct {
	// name is the short name or alias for the imported package. This is either the package's
	// default name, or the alias specified in the import statement.
	// Not used if self is true.
	name string

	// path is the full path to the imported package.
	path string

	// self is true when this is an implicit import declaration for the current package.
	self bool
}

// importSpecToDecl returns an import declaration for spec.
func importSpecToDecl(spec *ast.ImportSpec, imports []*types.Package) (importDeclaration, bool) {
	path := strings.TrimSuffix(strings.TrimPrefix(spec.Path.Value, `"`), `"`)

	if spec.Name != nil {
		return importDeclaration{
			name: spec.Name.Name,
			path: path,
		}, true
	}

	for _, imp := range imports {
		if imp.Path() == path {
			return importDeclaration{
				name: imp.Name(),
				path: path,
			}, true
		}
	}

	return importDeclaration{}, false
}
