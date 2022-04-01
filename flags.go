package varnamelen

import "strings"

// stringsValue is the value of a list-of-strings flag.
// FIXME: type
type stringsValue struct {
	Values []string
}

// declarationsValue is the value of a list-of-declarations flag.
// FIXME: type
type declarationsValue struct {
	Values []identDeclaration
}

// Set implements Value.
func (sv *stringsValue) Set(values string) error {
	if strings.TrimSpace(values) == "" {
		sv.Values = nil
		return nil
	}

	parts := strings.Split(values, ",")

	sv.Values = make([]string, len(parts))

	for i, part := range parts {
		sv.Values[i] = strings.TrimSpace(part)
	}

	return nil
}

// String implements Value.
func (sv *stringsValue) String() string {
	return strings.Join(sv.Values, ",")
}

// contains returns true if sv contains s.
func (sv stringsValue) contains(s string) bool {
	for _, v := range sv.Values {
		if v == s {
			return true
		}
	}

	return false
}

// Set implements Value.
func (dv *declarationsValue) Set(values string) error {
	if strings.TrimSpace(values) == "" {
		dv.Values = nil
		return nil
	}

	parts := strings.Split(values, ",")

	dv.Values = make([]identDeclaration, len(parts))

	for idx, part := range parts {
		var err error
		if dv.Values[idx], err = parseIdentDeclaration(strings.TrimSpace(part)); err != nil {
			return err
		}
	}

	return nil
}

// String implements Value.
func (dv *declarationsValue) String() string {
	parts := make([]string, len(dv.Values))

	for idx, val := range dv.Values {
		parts[idx] = val.name + " " + val.typ
	}

	return strings.Join(parts, ",")
}

// matchVariable returns true if v matches any of the declarations in dv.
func (dv declarationsValue) matchVariable(v variable) bool {
	for _, decl := range dv.Values {
		if v.match(decl) {
			return true
		}
	}

	return false
}

// matchVariable returns true if c matches any of the declarations in dv.
func (dv declarationsValue) matchConstant(c constant) bool {
	for _, decl := range dv.Values {
		if c.match(decl) {
			return true
		}
	}

	return false
}

// matchParameter returns true if p matches any of the declarations in dv.
func (dv declarationsValue) matchParameter(p parameter) bool {
	for _, decl := range dv.Values {
		if p.match(decl) {
			return true
		}
	}

	return false
}

// matchParameter returns true if r matches any of the declarations in dv.
func (dv declarationsValue) matchNamedReturn(r namedReturn) bool {
	for _, decl := range dv.Values {
		if r.match(decl) {
			return true
		}
	}

	return false
}

// matchParameter returns true if r matches any of the declarations in dv.
func (dv declarationsValue) matchReceiver(r receiver) bool {
	for _, decl := range dv.Values {
		if r.match(decl) {
			return true
		}
	}

	return false
}

// matchParameter returns true if p matches any of the declarations in dv.
func (dv declarationsValue) matchTypeParameter(p typeParam) bool {
	for _, decl := range dv.Values {
		if p.match(decl) {
			return true
		}
	}

	return false
}
