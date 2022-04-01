//go:build go1.18
// +build go1.18

package varnamelen

import "testing"

func TestRun_Warnings_TypeParam(t *testing.T) {
	runTest(t, "warningstypeparam", map[string]string{
		"checkTypeParam": "true",
	})
}

func TestRun_NameLen_TypeParam(t *testing.T) {
	runTest(t, "namelentypeparam", map[string]string{
		"checkTypeParam": "true",
	})
}

func TestRun_Distance_TypeParam(t *testing.T) {
	runTest(t, "distancetypeparam", map[string]string{
		"checkTypeParam": "true",
	})
}

func TestRun_IgnoreNames_TypeParam(t *testing.T) {
	runTest(t, "ignorenamestypeparam", map[string]string{
		"ignoreNames":    "T",
		"checkTypeParam": "true",
	})
}

func TestRun_IgnoreDecls_TypeParam(t *testing.T) {
	runTest(t, "ignoredeclstypeparam", map[string]string{
		"ignoreDecls":    "T any",
		"checkTypeParam": "true",
	})
}
