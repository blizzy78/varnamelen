package main //nolint:revive // don't need docs here

import (
	"github.com/blizzy78/varnamelen"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(varnamelen.NewAnalyzer())
}
