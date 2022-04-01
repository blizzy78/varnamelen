package varnamelen

import (
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRun_Warnings(t *testing.T) {
	runTest(t, "warnings", nil)
}

func TestRun_Warnings_Return(t *testing.T) {
	runTest(t, "warningsreturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestRun_Warnings_Receiver(t *testing.T) {
	runTest(t, "warningsreceiver", map[string]string{
		"checkReceiver": "true",
	})
}

func TestRun_NameLen(t *testing.T) {
	runTest(t, "namelen", nil)
}

func TestRun_NameLen_Return(t *testing.T) {
	runTest(t, "namelenreturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestRun_NameLen_Receiver(t *testing.T) {
	runTest(t, "namelenreceiver", map[string]string{
		"checkReceiver": "true",
	})
}

func TestRun_Distance(t *testing.T) {
	runTest(t, "distance", nil)
}

func TestRun_Distance_Return(t *testing.T) {
	runTest(t, "distancereturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestRun_Distance_Receiver(t *testing.T) {
	runTest(t, "distancereceiver", map[string]string{
		"checkReceiver": "true",
	})
}

func TestRun_IgnoreNames(t *testing.T) {
	runTest(t, "ignorenames", map[string]string{
		"ignoreNames": "i, I",
	})
}

func TestRun_IgnoreNames_Return(t *testing.T) {
	runTest(t, "ignorenamesreturn", map[string]string{
		"ignoreNames": "i",
		"checkReturn": "true",
	})
}

func TestRun_IgnoreNames_Receiver(t *testing.T) {
	runTest(t, "ignorenamesreceiver", map[string]string{
		"ignoreNames":   "f",
		"checkReceiver": "true",
	})
}

func TestRun_IgnoreDecls(t *testing.T) {
	runTest(t, "ignoredecls", map[string]string{
		"ignoreDecls": "i int, i map[string]string, i *bytes.Buffer, i *strs.Builder, const I",
	})
}

func TestRun_IgnoreDecls_Return(t *testing.T) {
	runTest(t, "ignoredeclsreturn", map[string]string{
		"ignoreDecls": "i int",
		"checkReturn": "true",
	})
}

func TestRun_IgnoreDecls_Receiver(t *testing.T) {
	runTest(t, "ignoredeclsreceiver", map[string]string{
		"ignoreDecls":   "f foo, f *foo",
		"checkReceiver": "true",
	})
}

func TestRun_IgnoreTypeAssertOk(t *testing.T) {
	runTest(t, "typeassert", map[string]string{
		"ignoreTypeAssertOk": "true",
	})
}

func TestRun_IgnoreChanRecvOk(t *testing.T) {
	runTest(t, "chanrecv", map[string]string{
		"ignoreChanRecvOk": "true",
	})
}

func TestRun_IgnoreMapIndexOk(t *testing.T) {
	runTest(t, "mapindex", map[string]string{
		"ignoreMapIndexOk": "true",
	})
}

func TestRun_Conventional(t *testing.T) {
	runTest(t, "conventional", nil)
}

func TestRun_FalseNegative_TypeAssertOk(t *testing.T) {
	runTest(t, "falsenegativetypeassert", map[string]string{
		"ignoreTypeAssertOk": "true",
	})
}

func TestRun_FalseNegative_ChanRecvOk(t *testing.T) {
	runTest(t, "falsenegativechanrecv", map[string]string{
		"ignoreChanRecvOk": "true",
	})
}

func TestRun_FalseNegative_MapIndexOk(t *testing.T) {
	runTest(t, "falsenegativemapindex", map[string]string{
		"ignoreMapIndexOk": "true",
	})
}

func runTest(t *testing.T, pkg string, flags map[string]string) {
	t.Helper()

	analyzer := NewAnalyzer()

	for k, v := range flags {
		_ = analyzer.Flags.Set(k, v)
	}

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, pkg)
}
