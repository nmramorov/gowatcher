// Модуль exitchecker проверяет не было ли вызовов функции os.Exit()
// в функции "main" модуля "main" для каждого модуля проекта.
package exitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for unchecked os.Exit usages",
	Run:  run,
}

// Сначала проверяем, что пакет называется "main", далее нужно проверить что
// узел синтаксического дерева может быть приведен к типу SelectorExpr,
// это необходимо для проверки модуля и функции os.Exit.
func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if x, ok := node.(*ast.SelectorExpr); ok {
					if pkg, ok := x.X.(*ast.Ident); ok {
						if pkg.Name == "os" && x.Sel.Name == "Exit" {
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
