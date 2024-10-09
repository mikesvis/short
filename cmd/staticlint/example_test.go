// Статический анализатор
package main

import (
	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func ExampleMain() {
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
