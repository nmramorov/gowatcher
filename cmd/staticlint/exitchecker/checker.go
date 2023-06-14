package exitchecker

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for unchecked os.Exit usages",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if strings.Compare(pass.Pkg.Name(), "main") == 0 {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if x, ok := node.(*ast.SelectorExpr); ok {
					if pkg, ok := x.X.(*ast.Ident); ok {
						if strings.Compare(pkg.Name, "os") == 0 && strings.Compare(x.Sel.Name, "Exit") == 0 {
							pass.Reportf(x.Pos(), "os.Exit used in main!!!")
						}
					}

				}
				return true
			})
		}

	}
	return nil, nil
}
