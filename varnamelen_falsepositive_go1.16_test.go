package varnamelen

import "testing"

func TestRun_FalsePositive_Sixteen(t *testing.T) {
	runTest(t, "falsepositivesixteen", nil)
}
