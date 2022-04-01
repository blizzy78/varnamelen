//go:build go1.18
// +build go1.18

package varnamelen

import "testing"

func TestRun_FalsePositive(t *testing.T) {
	runTest(t, "falsepositive", nil)
}
