package main

import (
	"go/ast"
	"log"

	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

// Статический анализатор для проекта. Запуск:
//
//	go vet -vettool=./cmd/staticlint/staticlint ./...
func main() {
	var analyzers []*analysis.Analyzer

	// Стандартные анализаторы golang.org/x/tools/go/analysis/passes
	analyzers = append(analyzers, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)

	// Анализаторы класса SA пакета staticcheck.io
	analyzers = append(analyzers, getStaticCheckAnalyzers()...)

	// Публичный анализатор (в нем несколько анализаторов)
	analyzers = append(analyzers, analyzer.Analyzer)

	// Кастомный анализатор проверки os.Exit
	analyzers = append(analyzers, ExitAnalyzer)

	multichecker.Main(analyzers...)
}

func getStaticCheckAnalyzers() []*analysis.Analyzer {
	var result []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		result = append(result, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		result = append(result, v.Analyzer)
	}

	return result
}

// Анализатор использования os.Exit() в ф-ции main
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "Check os.Exit call in main function",
	Run:  run,
}

func isOsExitCall(callExpr *ast.CallExpr) bool {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if selExpr.Sel.Name != "Exit" {
		return false
	}
	pkgIdent, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkgIdent.Name == "os"
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		var inMainFunc bool

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				inMainFunc = x.Name.Name == "main"
			case *ast.CallExpr:
				if inMainFunc && isOsExitCall(x) {
					log.Println("Found a call to os.Exit in main function")
					pass.Reportf(x.Pos(), "call to os.Exit in main function is voilated")
				}
			}
			return true
		})
	}
	return nil, nil
}
