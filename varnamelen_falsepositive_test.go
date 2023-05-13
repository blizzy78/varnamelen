package varnamelen

import "testing"

func TestRun_FalsePositive(t *testing.T) {
	runTest(t, "falsepositive", nil)
}
