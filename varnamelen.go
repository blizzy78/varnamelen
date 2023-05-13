package varnamelen

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// configuration is the configuration of this analyzer.
type configuration struct {
	// maxDistance is the longest distance, in source lines, that is being considered a "small" scope.
	maxDistance int

	// minNameLength is the minimum length of a name that is considered "long."
	minNameLength int

	// ignoreNames is an optional list of names that should be ignored completely.
	ignoreNames stringsValue

	// checkReceiver determines whether method receivers should be checked.
	checkReceiver bool

	// checkReturn determines whether named return values should be checked.
	checkReturn bool

	// checkTypeParam determines whether type parameters should be checked.
	checkTypeParam bool

	// ignoreTypeAssertOk determines whether "ok" variables that hold the bool return value of a type assertion should be ignored.
	ignoreTypeAssertOk bool

	// ignoreMapIndexOk determines whether "ok" variables that hold the bool return value of a map index should be ignored.
	ignoreMapIndexOk bool

	// ignoreChanRecvOk determines whether "ok" variables that hold the bool return value of a channel receive should be ignored.
	ignoreChanRecvOk bool

	// ignoreDecls is an optional list of declarations that should be ignored completely.
	ignoreDecls declarationsValue
}

// identsImportsSwitchesResult is the result of calling the identsImportsSwitches function.
type identsImportsSwitchesResult struct {
	// assignIdents is a list of idents used in assign statements.
	assignIdents []*ast.Ident

	// valueSpecIdents is a list of idents used in value spec (var) statements.
	valueSpecIdents []*ast.Ident

	// paramIdents is a list of idents used in parameters.
	paramIdents []*ast.Ident

	// returnIdents is a list of idents used in named return values.
	returnIdents []*ast.Ident

	// receiverIdents is a list of idents used in method receivers.
	receiverIdents []*ast.Ident

	// typeParamIdents is a list of idents used in type parameters.
	typeParamIdents []*ast.Ident

	// imports is a list of import declarations.
	imports []importDeclaration

	// switches is a list of type switch statements.
	switches []*ast.TypeSwitchStmt
}

// distancesResult is the result of calling the distances function.
type distancesResult struct {
	// variableToDist is a map of variables to their longest usage distance.
	variableToDist map[variable]int

	// constantToDist is a map of constants to their longest usage distance.
	constantToDist map[constant]int

	// paramToDist is a map of parameters to their longest usage distance.
	paramToDist map[parameter]int

	// returnToDist is a map of named return values to their longest usage distance.
	returnToDist map[namedReturn]int

	// receiverToDist is a map of method receivers to their longest usage distance.
	receiverToDist map[receiver]int

	// typeParamToDist is a map of type parameters to their longest usage distance.
	typeParamToDist map[typeParam]int
}

const (
	// defaultMaxDistance is the default value for the maximum distance between the declaration of an ident and its usage
	// that is considered a "small" scope.
	defaultMaxDistance = 5

	// defaultMinNameLength is the default value for the minimum length of an ident's name that is considered "long."
	defaultMinNameLength = 3
)

// conventionalDecls is a list of conventional declarations. These will always be ignored.
var conventionalDecls = []identDeclaration{
	mustParseIdentDeclaration("ctx context.Context"),

	mustParseIdentDeclaration("b *testing.B"),
	mustParseIdentDeclaration("f *testing.F"),
	mustParseIdentDeclaration("m *testing.M"),
	mustParseIdentDeclaration("pb *testing.PB"),
	mustParseIdentDeclaration("t *testing.T"),
	mustParseIdentDeclaration("tb testing.TB"),
}

// NewAnalyzer returns a new analyzer.
func NewAnalyzer() *analysis.Analyzer {
	cfg := configuration{
		maxDistance:   defaultMaxDistance,
		minNameLength: defaultMinNameLength,
		ignoreNames:   stringsValue{},
		ignoreDecls:   declarationsValue{},
	}

	analyzer := analysis.Analyzer{
		Name: "varnamelen",
		Doc: "checks that the length of a variable's name matches its scope\n\n" +
			"Variables with short names can be hard to use if the variable is used\n" +
			"over a longer span of lines of code. A longer variable name may be easier\n" +
			"to comprehend. Also checks constants, parameters, named return values,\n" +
			"method receivers, and type parameters",

		Run: func(pass *analysis.Pass) (any, error) {
			run(pass, cfg)
			return nil, nil
		},

		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	analyzer.Flags.IntVar(&cfg.maxDistance, "maxDistance", cfg.maxDistance, "maximum number of lines considered a 'short' scope")
	analyzer.Flags.IntVar(&cfg.minNameLength, "minNameLength", cfg.minNameLength, "minimum length of a name considered 'long'")
	analyzer.Flags.Var(&cfg.ignoreNames, "ignoreNames", "comma-separated list of ignored names")
	analyzer.Flags.BoolVar(&cfg.checkReceiver, "checkReceiver", cfg.checkReceiver, "check method receivers")
	analyzer.Flags.BoolVar(&cfg.checkReturn, "checkReturn", cfg.checkReturn, "check named return values")
	analyzer.Flags.BoolVar(&cfg.ignoreTypeAssertOk, "ignoreTypeAssertOk", cfg.ignoreTypeAssertOk, "ignore 'ok' variables that hold the bool return value of a type assertion")
	analyzer.Flags.BoolVar(&cfg.ignoreMapIndexOk, "ignoreMapIndexOk", cfg.ignoreMapIndexOk, "ignore 'ok' variables that hold the bool return value of a map index")
	analyzer.Flags.BoolVar(&cfg.ignoreChanRecvOk, "ignoreChanRecvOk", cfg.ignoreChanRecvOk, "ignore 'ok' variables that hold the bool return value of a channel receive")
	analyzer.Flags.Var(&cfg.ignoreDecls, "ignoreDecls", "comma-separated list of ignored declarations")
	analyzer.Flags.BoolVar(&cfg.checkTypeParam, "checkTypeParam", cfg.checkTypeParam, "check type parameters")

	return &analyzer
}

// run applies the analysis to pass, according to cfg.
func run(pass *analysis.Pass, cfg configuration) {
	dist := distances(pass, cfg)

	checkVariables(pass, dist.variableToDist, cfg)
	checkConstants(pass, dist.constantToDist, cfg)
	checkParams(pass, dist.paramToDist, cfg)
	checkReturns(pass, dist.returnToDist, cfg)
	checkReceivers(pass, dist.receiverToDist, cfg)
	checkTypeParams(pass, dist.typeParamToDist, cfg)
}

// distances returns maps of idents to their longest usage distance.
func distances(pass *analysis.Pass, cfg configuration) distancesResult {
	result := distancesResult{
		variableToDist:  map[variable]int{},
		constantToDist:  map[constant]int{},
		paramToDist:     map[parameter]int{},
		returnToDist:    map[namedReturn]int{},
		receiverToDist:  map[receiver]int{},
		typeParamToDist: map[typeParam]int{},
	}

	idents := identsImportsSwitches(pass, cfg)

	for _, ident := range idents.assignIdents {
		assign := ident.Obj.Decl.(*ast.AssignStmt) //nolint:forcetypeassert // check is done in identsAndImports

		var typ string
		if isTypeSwitchAssign(assign, idents.switches) {
			typ = "<type-switched>"
		} else {
			typ = shortTypeName(pass.TypesInfo.TypeOf(ident), idents.imports)
		}

		variable := variable{
			name:   ident.Name,
			typ:    typ,
			assign: assign,
		}

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(assign.Pos()).Line
		result.variableToDist[variable] = useLine - declLine
	}

	for _, ident := range idents.valueSpecIdents {
		valueSpec := ident.Obj.Decl.(*ast.ValueSpec) //nolint:forcetypeassert // check is done in identsAndImports

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(valueSpec.Pos()).Line

		if ident.Obj.Kind == ast.Con {
			constant := constant{
				name:      ident.Name,
				valueSpec: valueSpec,
			}

			result.constantToDist[constant] = useLine - declLine
		} else {
			variable := variable{
				name:      ident.Name,
				typ:       shortTypeName(pass.TypesInfo.TypeOf(ident), idents.imports),
				valueSpec: valueSpec,
			}

			result.variableToDist[variable] = useLine - declLine
		}
	}

	for _, ident := range idents.paramIdents {
		field := ident.Obj.Decl.(*ast.Field) //nolint:forcetypeassert // check is done in identsAndImports

		param := parameter{
			name:  ident.Name,
			typ:   shortTypeName(pass.TypesInfo.TypeOf(field.Type), idents.imports),
			field: field,
		}

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(field.Pos()).Line
		result.paramToDist[param] = useLine - declLine
	}

	for _, ident := range idents.returnIdents {
		field := ident.Obj.Decl.(*ast.Field) //nolint:forcetypeassert // check is done in identsAndImports

		param := namedReturn{
			name:  ident.Name,
			typ:   shortTypeName(pass.TypesInfo.TypeOf(ident), idents.imports),
			field: field,
		}

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(field.Pos()).Line
		result.returnToDist[param] = useLine - declLine
	}

	for _, ident := range idents.receiverIdents {
		field := ident.Obj.Decl.(*ast.Field) //nolint:forcetypeassert // check is done in identsAndImports

		param := receiver{
			name:  ident.Name,
			typ:   shortTypeName(pass.TypesInfo.TypeOf(field.Type), idents.imports),
			field: field,
		}

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(field.Pos()).Line
		result.receiverToDist[param] = useLine - declLine
	}

	for _, ident := range idents.typeParamIdents {
		field := ident.Obj.Decl.(*ast.Field) //nolint:forcetypeassert // check is done in identsAndImports

		param := typeParam{
			name:  ident.Name,
			typ:   shortTypeName(pass.TypesInfo.TypeOf(field.Type), idents.imports),
			field: field,
		}

		useLine := pass.Fset.Position(ident.NamePos).Line
		declLine := pass.Fset.Position(field.Pos()).Line
		result.typeParamToDist[param] = useLine - declLine
	}

	return result
}

// identsImportsSwitches returns relevant idents, as well as import declarations and type switch statements.
func identsImportsSwitches(pass *analysis.Pass, cfg configuration) identsImportsSwitchesResult { //nolint:gocognit,cyclop // this is complex stuff
	result := identsImportsSwitchesResult{}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector) //nolint:forcetypeassert // inspect.Analyzer always returns *inspector.Inspector

	filter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.TypeSwitchStmt)(nil),
		(*ast.Ident)(nil),
	}

	funcs := []*ast.FuncDecl{}
	methods := []*ast.FuncDecl{}
	funcLits := []*ast.FuncLit{}
	compositeLits := []*ast.CompositeLit{}

	inspector.Preorder(filter, func(node ast.Node) {
		switch node2 := node.(type) {
		case *ast.ImportSpec:
			decl, ok := importSpecToDecl(node2, pass.Pkg.Imports())
			if !ok {
				return
			}

			result.imports = append(result.imports, decl)

		case *ast.FuncDecl:
			funcs = append(funcs, node2)

			if node2.Recv == nil {
				return
			}

			methods = append(methods, node2)

		case *ast.FuncLit:
			funcLits = append(funcLits, node2)

		case *ast.CompositeLit:
			compositeLits = append(compositeLits, node2)

		case *ast.TypeSwitchStmt:
			result.switches = append(result.switches, node2)

		case *ast.Ident:
			if node2.Obj == nil {
				return
			}

			if isCompositeLitKey(node2, compositeLits) {
				return
			}

			switch objDecl := node2.Obj.Decl.(type) {
			case *ast.AssignStmt:
				result.assignIdents = append(result.assignIdents, node2)

			case *ast.ValueSpec:
				result.valueSpecIdents = append(result.valueSpecIdents, node2)

			case *ast.Field:
				switch {
				case isParam(objDecl, funcs, funcLits, methods):
					result.paramIdents = append(result.paramIdents, node2)

				case isReturn(objDecl, funcs, funcLits):
					if !cfg.checkReturn {
						return
					}

					result.returnIdents = append(result.returnIdents, node2)

				case isReceiver(objDecl, methods):
					if !cfg.checkReceiver {
						return
					}

					result.receiverIdents = append(result.receiverIdents, node2)

				case isTypeParam(objDecl, funcs, funcLits):
					if !cfg.checkTypeParam {
						return
					}

					result.typeParamIdents = append(result.typeParamIdents, node2)
				}
			}
		}
	})

	result.imports = append(result.imports, importDeclaration{
		path: pass.Pkg.Path(),
		self: true,
	})

	sort.Slice(result.imports, func(a int, b int) bool {
		// reversed: longest path first
		return len(result.imports[a].path) > len(result.imports[b].path)
	})

	return result
}

// isReceiver returns true if field is the receiver of any of the given methods.
func isReceiver(field *ast.Field, methods []*ast.FuncDecl) bool {
	for _, m := range methods {
		for _, recv := range m.Recv.List {
			if recv == field {
				return true
			}
		}
	}

	return false
}

// isReturn returns true if field is a return value of any of the given funcs.
func isReturn(field *ast.Field, funcs []*ast.FuncDecl, funcLits []*ast.FuncLit) bool { //nolint:gocognit // it's not that complicated
	for _, f := range funcs {
		if f.Type.Results == nil {
			continue
		}

		for _, r := range f.Type.Results.List {
			if r == field {
				return true
			}
		}
	}

	for _, f := range funcLits {
		if f.Type.Results == nil {
			continue
		}

		for _, r := range f.Type.Results.List {
			if r == field {
				return true
			}
		}
	}

	return false
}

// isParam returns true if field is a parameter of any of the given funcs or methods.
func isParam(field *ast.Field, funcs []*ast.FuncDecl, funcLits []*ast.FuncLit, methods []*ast.FuncDecl) bool { //nolint:gocognit,cyclop // it's not that complicated
	for _, f := range funcs {
		if f.Type.Params == nil {
			continue
		}

		for _, p := range f.Type.Params.List {
			if p == field {
				return true
			}
		}
	}

	for _, f := range funcLits {
		if f.Type.Params == nil {
			continue
		}

		for _, p := range f.Type.Params.List {
			if p == field {
				return true
			}
		}
	}

	for _, m := range methods {
		if m.Type.Params == nil {
			continue
		}

		for _, p := range m.Type.Params.List {
			if p == field {
				return true
			}
		}
	}

	return false
}

// isCompositeLitKey returns true if ident is a key of any of the given composite literals.
func isCompositeLitKey(ident *ast.Ident, compositeLits []*ast.CompositeLit) bool {
	for _, cl := range compositeLits {
		if _, ok := cl.Type.(*ast.MapType); ok {
			continue
		}

		for _, kvExpr := range cl.Elts {
			kv, ok := kvExpr.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			if kv.Key == ident {
				return true
			}
		}
	}

	return false
}

// isTypeSwitchAssign returns true if assign is an assign statement of any of the given type switch statements.
func isTypeSwitchAssign(assign *ast.AssignStmt, switches []*ast.TypeSwitchStmt) bool {
	for _, s := range switches {
		if s.Assign == assign {
			return true
		}
	}

	return false
}

// shortTypeName returns the short name of typ, with respect to imports.
// For example, if package github.com/matryer/is is imported with alias "x",
// and typ represents []*github.com/matryer/is.I, shortTypeName will return "[]*x.I".
// For imports without aliases, the package's default name will be used.
func shortTypeName(typ types.Type, imports []importDeclaration) string {
	if typ == nil {
		return ""
	}

	typStr := typ.String()

	for _, imp := range imports {
		prefix := imp.path + "."

		replace := ""
		if !imp.self {
			replace = imp.name + "."
		}

		typStr = strings.ReplaceAll(typStr, prefix, replace)
	}

	return typStr
}

// checkNameAndDistance returns true if the named ident should be ignored, because its name
// is "long" enough or its usage distance is to be considered "short", according to cfg.
func checkNameAndDistance(name string, dist int, cfg configuration) bool {
	if len([]rune(name)) >= cfg.minNameLength {
		return true
	}

	if dist <= cfg.maxDistance {
		return true
	}

	return false
}
