package varnamelen

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarNameLen_Run(t *testing.T) {
	a := NewAnalyzer()
	a.Flags.Set("ignoreNames", "i")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "test")
}

func TestIntValue_Set(t *testing.T) {
	is := is.New(t)
	v := intValue(0)
	v.Set("123")
	is.Equal(int(v), 123)
}

func TestIntValue_Set_Error(t *testing.T) {
	is := is.New(t)
	v := intValue(0)
	is.True(v.Set("") != nil)
	is.True(v.Set("test") != nil)
}

func TestIntValue_String(t *testing.T) {
	is := is.New(t)
	v := intValue(123)
	is.Equal(v.String(), "123")
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
