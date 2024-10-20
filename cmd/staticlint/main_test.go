// Набор статических анализаторов для проекта.
package main

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

func Test_getStaticCheckAnalyzers(t *testing.T) {
	tests := []struct {
		name string
		want []*analysis.Analyzer
	}{
		{
			name: "Has type of",
			want: []*analysis.Analyzer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStaticCheckAnalyzers()
			assert.IsType(t, tt.want, got)
		})
	}
}

func Test_isOsExitCall(t *testing.T) {
	type args struct {
		callExpr *ast.CallExpr
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Expression",
			args: args{
				callExpr: &ast.CallExpr{},
			},
			want: false,
		},
		{
			name: "Not an exit expr",
			args: args{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "NOPE",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Not an exit expr",
			args: args{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "NOPE",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Not an os.Exit expr",
			args: args{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "Exit",
						},
						X: &ast.Ident{
							Name: "win",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Is os.Exit expr",
			args: args{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "Exit",
						},
						X: &ast.Ident{
							Name: "os",
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOsExitCall(tt.args.callExpr)
			assert.Equal(t, tt.want, got)
		})
	}
}
