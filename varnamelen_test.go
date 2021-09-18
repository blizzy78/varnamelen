package varnamelen

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarNameLen_Run(t *testing.T) {
	a := NewAnalyzer()
	a.Flags.Set("minNameLength", "4")
	a.Flags.Set("ignoreNames", "i,ip")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "test")
}

func TestVarNameLen_Run_CheckReceiver(t *testing.T) {
	a := NewAnalyzer()
	a.Flags.Set("minNameLength", "4")
	a.Flags.Set("checkReceiver", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "receiver")
}

func TestVarNameLen_Run_CheckReturn(t *testing.T) {
	a := NewAnalyzer()
	a.Flags.Set("minNameLength", "4")
	a.Flags.Set("ignoreNames", "i")
	a.Flags.Set("checkReturn", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "return")
}

func TestStringsValue_Set(t *testing.T) {
	is := is.New(t)
	v := stringsValue{}
	v.Set("foo,bar,baz")
	is.Equal(v.Values, []string{"foo", "bar", "baz"})
}

func TestStringsValue_String(t *testing.T) {
	is := is.New(t)
	v := stringsValue{
		Values: []string{"foo", "bar", "baz"},
	}
	is.Equal(v.String(), "foo,bar,baz")
}
