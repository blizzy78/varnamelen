package varnamelen

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarNameLen_Run(t *testing.T) {
	a := NewAnalyzer()
	_ = a.Flags.Set("minNameLength", "4")
	_ = a.Flags.Set("ignoreNames", "i, ip, CI")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "test")
}

func TestVarNameLen_Run_CheckReceiver(t *testing.T) {
	a := NewAnalyzer()
	_ = a.Flags.Set("minNameLength", "4")
	_ = a.Flags.Set("checkReceiver", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", a, "receiver")
}

func TestVarNameLen_Run_CheckReturn(t *testing.T) {
	analyzer := NewAnalyzer()
	_ = analyzer.Flags.Set("minNameLength", "4")
	_ = analyzer.Flags.Set("ignoreNames", "i")
	_ = analyzer.Flags.Set("checkReturn", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, "return")
}

func TestVarNameLen_Run_IgnoreTypeAssertOk(t *testing.T) {
	analyzer := NewAnalyzer()
	_ = analyzer.Flags.Set("minNameLength", "4")
	_ = analyzer.Flags.Set("ignoreTypeAssertOk", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, "type-assert-ok")
}

func TestVarNameLen_Run_IgnoreMapIndexOk(t *testing.T) {
	analyzer := NewAnalyzer()
	_ = analyzer.Flags.Set("minNameLength", "4")
	_ = analyzer.Flags.Set("ignoreMapIndexOk", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, "map-index-ok")
}

func TestVarNameLen_Run_IgnoreChannelReceiveOk(t *testing.T) {
	analyzer := NewAnalyzer()
	_ = analyzer.Flags.Set("minNameLength", "4")
	_ = analyzer.Flags.Set("ignoreChanRecvOk", "true")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, "chan-recv-ok")
}

func TestVarNameLen_Run_IgnoreDeclarations(t *testing.T) {
	analyzer := NewAnalyzer()
	_ = analyzer.Flags.Set("minNameLength", "4")
	_ = analyzer.Flags.Set("checkReturn", "true")
	_ = analyzer.Flags.Set("ignoreDecls", "c context.Context, b bb.Buffer, b *strings.Builder, d *bb.Buffer, i int, ip *int, const C, f func(), m map[int]*bb.Buffer, mi int")

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, "decl")
}

func TestStringsValue_Set(t *testing.T) {
	is := is.New(t)
	v := stringsValue{}
	_ = v.Set("foo,bar,baz")
	is.Equal(v.Values, []string{"foo", "bar", "baz"})
}

func TestStringsValue_String(t *testing.T) {
	is := is.New(t)
	v := stringsValue{
		Values: []string{"foo", "bar", "baz"},
	}
	is.Equal(v.String(), "foo,bar,baz")
}

func TestParseDeclaration(t *testing.T) {
	tests := []struct {
		givenDecl string
		wantDecl  declaration
	}{
		{
			givenDecl: "t *testing.T",
			wantDecl:  declaration{name: "t", typ: "*testing.T"},
		},
		{
			givenDecl: "c echo.Context",
			wantDecl:  declaration{name: "c", typ: "echo.Context"},
		},
		{
			givenDecl: "const C",
			wantDecl:  declaration{name: "C", constant: true},
		},
	}

	for _, test := range tests {
		t.Run(test.givenDecl, func(t *testing.T) {
			is := is.New(t)

			is.Equal(parseDeclaration(test.givenDecl), test.wantDecl)
		})
	}
}
