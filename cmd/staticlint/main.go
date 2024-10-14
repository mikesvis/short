// Набор статических анализаторов для проекта.
//
// Запуск:
//
//	go vet -vettool=./cmd/staticlint/staticlint ./...
//
// Набор анализаторов:
//
// Стандартные (golang.org/x/tools/go/analysis/passes):
//
//	printf: check consistency of Printf format strings and arguments
//	shadow: check for possible unintended shadowing of variables
//	structtag defines an Analyzer that checks struct field tags are well formed.
//
// Анализаторы класса SA пакета staticcheck.io:
//
//	sa1000 Invalid regular expression
//	sa1001 Invalid template
//	sa1002 Invalid format in 'time.Parse'
//	sa1003 Unsupported argument to functions in 'encoding/binary'
//	sa1004 Suspiciously small untyped constant in 'time.Sleep'
//	sa1005 Invalid first argument to 'exec.Command'
//	sa1006 'Printf' with dynamic first argument and no further arguments
//	sa1007 Invalid URL in 'net/url.Parse'
//	sa1008 Non-canonical key in 'http.Header' map
//	sa1010 '(*regexp.Regexp).FindAll' called with 'n == 0', which will always return zero results
//	sa1011 `Various methods in the "strings" package expect valid UTF-8, but invalid input is provided
//	sa1012 A nil 'context.Context' is being passed to a function, consider using 'context.TODO' instead
//	sa1013 'io.Seeker.Seek' is being called with the whence constant as the first argument, but it should be the second
//	sa1014 Non-pointer value passed to 'Unmarshal' or 'Decode\
//	sa1015 Using 'time.Tick' in a way that will leak. Consider using 'time.NewTicker', and only use 'time.Tick' in tests, commands and endless functions
//	sa1016 Trapping a signal that cannot be trapped
//	sa1017 Channels used with 'os/signal.Notify' should be buffered
//	sa1018 'strings.Replace' called with 'n == 0', which does nothing
//	sa1019 Using a deprecated function, variable, constant or field
//	sa1020 Using an invalid host:port pair with a 'net.Listen'-related function
//	sa1021 Using 'bytes.Equal' to compare two 'net.IP'
//	sa1023 Modifying the buffer in an 'io.Writer' implementation
//	sa1024 A string cutset contains duplicate characters
//	sa1025 It is not possible to use '(*time.Timer).Reset''s return value correctly
//	sa1026 Cannot marshal channels or functions
//	sa1027 Atomic access to 64-bit variable must be 64-bit aligned
//	sa1028 'sort.Slice' can only be used on slices
//	sa1029 Inappropriate key in call to 'context.WithValue'
//	sa1030 Invalid argument in call to a 'strconv' function
//	sa1031 Overlapping byte slices passed to an encoder
//	sa1032 Wrong order of arguments to 'errors.Is'
//	sa2000 'sync.WaitGroup.Add' called inside the goroutine, leading to a race condition
//	sa2001 Empty critical section, did you mean to defer the unlock?
//	sa2002 Called 'testing.T.FailNow' or 'SkipNow' in a goroutine, which isn't allowed
//	sa2003 Deferred 'Lock' right after locking, likely meant to defer 'Unlock' instead
//	sa3000 'TestMain' doesn't call 'os.Exit', hiding test failures
//	sa3001 Assigning to 'b.N' in benchmarks distorts the results
//	sa4000 Binary operator has identical expressions on both sides
//	sa4001 '&*x' gets simplified to 'x', it does not copy 'x'
//	sa4003 Comparing unsigned values against negative values is pointless
//	sa4004 The loop exits unconditionally after one iteration
//	sa4005 Field assignment that will never be observed. Did you mean to use a pointer receiver?
//	sa4006 A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
//	sa4008 The variable in the loop condition never changes, are you incrementing the wrong variable?
//	sa4009 A function argument is overwritten before its first use
//	sa4010 The result of 'append' will never be observed anywhere
//	sa4011 Break statement with no effect. Did you mean to break out of an outer loop?
//	sa4012 Comparing a value against NaN even though no value is equal to NaN
//	sa4013 Negating a boolean twice ('!!b') is the same as writing 'b'. This is either redundant, or a typo.
//	sa4014 An if/else if chain has repeated conditions and no side-effects; if the condition didn't match the first time, it won't match the second time, either
//	sa4015 Calling functions like 'math.Ceil' on floats converted from integers doesn't do anything useful
//	sa4016 Certain bitwise operations, such as 'x ^ 0', do not do anything useful
//	sa4017 Discarding the return values of a function without side effects, making the call pointless
//	sa4018 Self-assignment of variables
//	sa4019 Multiple, identical build constraints in the same file
//	sa4020 Unreachable case clause in a type switch
//	sa4021 "x = append(y)" is equivalent to "x = y"
//	sa4022 Comparing the address of a variable against nil
//	sa4023 Impossible comparison of interface value with untyped nil
//	sa4024 Checking for impossible return value from a builtin function
//	sa4025 Integer division of literals that results in zero
//	sa4026 Go constants cannot express negative zero
//	sa4027 '(*net/url.URL).Query' returns a copy, modifying it doesn't change the URL
//	sa4028 'x % 1' is always zero
//	sa4029 Ineffective attempt at sorting slice
//	sa4030 Ineffective attempt at generating random number
//	sa4031 Checking never-nil value against nil
//	sa4032 Comparing 'runtime.GOOS' or 'runtime.GOARCH' against impossible value
//	sa5000 Assignment to nil map
//	sa5001 Deferring 'Close' before checking for a possible error
//	sa5002 The empty for loop (\"for {}\") spins and can block the scheduler
//	sa5003 Defers in infinite loops will never execute
//	sa5004 "for { select { ..." with an empty default branch spins
//	sa5005 The finalizer references the finalized object, preventing garbage collection
//	sa5007 Infinite recursive call
//	sa5008 Invalid struct tag
//	sa5009 Invalid Printf call
//	sa5010 Impossible type assertion
//	sa5011 Possible nil pointer dereference
//	sa5012 Passing odd-sized slice to function expecting even size
//	sa6000 Using 'regexp.Match' or related in a loop, should use 'regexp.Compile'
//	sa6001 Missing an optimization opportunity when indexing maps by byte slices
//	sa6002 Storing non-pointer values in 'sync.Pool' allocates memory
//	sa6003 Converting a string to a slice of runes before ranging over it
//	sa6005 Inefficient string comparison with 'strings.ToLower' or 'strings.ToUpper'
//	sa6006 Using io.WriteString to write '[]byte'
//	sa9001 Defers in range loops may not run when you expect them to
//	sa9002 Using a non-octal 'os.FileMode' that looks like it was meant to be in octal.
//	sa9003 Empty body in an if or else branch
//	sa9004 Only the first constant has an explicit type
//	sa9005 Trying to marshal a struct with no public fields nor custom marshaling
//	sa9006 Dubious bit shifting of a fixed size integer value
//	sa9007 Deleting a directory that shouldn't be deleted
//	sa9008 'else' branch of a type assertion is probably not reading the right value
//	sa9009 Ineffectual Go compiler directive
//
// Кастомные анализаторы:
//
//	exitcheck Check os.Exit call in main function
//
// Сторонние анализаторы:
//
// go-critic Семейство аналитзаторов, включенных по умолчанию. Подробнее https://go-critic.com/overview
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
