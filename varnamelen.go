package varnamelen

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// varNameLen is an analyzer that checks that the length of a variable's name matches its usage scope.
// It will create a report for a variable's assignment if that variable has a short name, but its
// usage scope is not considered "small."
type varNameLen struct {
	// maxDistance is the longest distance, in source lines, that is being considered a "small scope."
	maxDistance intValue

	// minNameLength is the minimum length of a variable's name that is considered "long."
	minNameLength intValue

	// ignoreNames is an optional list of variable names that should be ignored completely.
	ignoreNames stringsValue
}

// intValue is the value of an integer flag.
type intValue int

// stringsValue is the value of a list-of-strings flag.
type stringsValue struct {
	Values []string
}

const (
	// defaultMaxDistance is the default value for the maximum distance between the declaration of a variable and its usage
	// that is considered a "small scope."
	defaultMaxDistance = 5

	// defaultMinNameLength is the default value for the minimum length of a variable's name that is considered "long."
	defaultMinNameLength = 3
)

// NewAnalyzer returns a new analyzer that uses varNameLen.
func NewAnalyzer() *analysis.Analyzer {
	vnl := varNameLen{
		maxDistance:   defaultMaxDistance,
		minNameLength: defaultMinNameLength,
		ignoreNames:   stringsValue{},
	}

	analyzer := analysis.Analyzer{
		Name: "varnamelen",
		Doc: "checks that the length of a variable's name matches its scope\n\n" +
			"A variable with a short name can be hard to use if the variable is used\n" +
			"over a longer span of lines of code. A longer variable name may be easier\n" +
			"to comprehend.",

		Run: func(pass *analysis.Pass) (interface{}, error) {
			vnl.run(pass)
			return nil, nil
		},

		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	analyzer.Flags.Var(&vnl.maxDistance, "maxDistance", "maximum number of lines of variable usage scope considered 'short'")
	analyzer.Flags.Var(&vnl.maxDistance, "minNameLength", "minimum length of variable name considered 'long'")
	analyzer.Flags.Var(&vnl.ignoreNames, "ignoreNames", "comma-separated list of ignored variable names")

	return &analyzer
}

// Run applies v to a package, according to pass.
func (v *varNameLen) run(pass *analysis.Pass) {
	nameToAssign, nameToDist := assignsAndDistances(pass)
	for name, assign := range nameToAssign {
		if len(name) >= int(v.minNameLength) {
			continue
		}
		if v.ignoreNames.contains(name) {
			continue
		}
		if nameToDist[name] <= int(v.maxDistance) {
			continue
		}

		pass.Reportf(assign.Pos(), "variable name '%s' is too short for the scope of its usage", name)
	}
}

// assignsAndDistances returns a map of variable names to assign statements, and a map of
// variable names to the longest usage distances. Both maps will have the same keys.
func assignsAndDistances(pass *analysis.Pass) (map[string]*ast.AssignStmt, map[string]int) {
	idents := identsReferencingAssigns(pass)

	nameToAssign := map[string]*ast.AssignStmt{}
	nameToDist := map[string]int{}

	for _, ident := range idents {
		useLine := pass.Fset.Position(ident.NamePos).Line
		assign := ident.Obj.Decl.(*ast.AssignStmt)
		declLine := pass.Fset.Position(assign.Pos()).Line

		dist := useLine - declLine
		if dist <= nameToDist[ident.Name] {
			continue
		}

		nameToAssign[ident.Name] = assign
		nameToDist[ident.Name] = dist
	}

	return nameToAssign, nameToDist
}

// identsReferencingAssigns returns all Idents in pass that reference assign statements.
func identsReferencingAssigns(pass *analysis.Pass) []*ast.Ident {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	filter := []ast.Node{
		(*ast.Ident)(nil),
	}

	idents := []*ast.Ident{}

	inspector.Preorder(filter, func(n ast.Node) {
		ident := n.(*ast.Ident)
		if ident.Obj == nil {
			return
		}
		if _, ok := ident.Obj.Decl.(*ast.AssignStmt); !ok {
			return
		}

		idents = append(idents, ident)
	})

	return idents
}

// Set implements Value.
func (i *intValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = intValue(v)
	return nil
}

// String implements Value.
func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

// Set implements Value.
func (sv *stringsValue) Set(s string) error {
	sv.Values = strings.Split(s, ",")
	return nil
}

// String implements Value.
func (sv *stringsValue) String() string {
	return strings.Join(sv.Values, ",")
}

func (sv *stringsValue) contains(s string) bool {
	for _, v := range sv.Values {
		if v == s {
			return true
		}
	}
	return false
}
